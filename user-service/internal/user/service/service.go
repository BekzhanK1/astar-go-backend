package user

import (
	"context"
	"errors"
	"fmt"
	userModels "user-service/internal/user/models"
	userRepo "user-service/internal/user/repository"
	"user-service/internal/utils"
)

type Service interface {
	Register(ctx context.Context, user *userModels.User) (*userModels.User, error)
	GetProfile(ctx context.Context, id uint) (*userModels.User, error)
	UpdateProfile(ctx context.Context, user *userModels.User) error
	DeleteUser(ctx context.Context, id uint) error
	ValidateUser(ctx context.Context, email, password string) (bool, *userModels.User, error)
}

type service struct {
	repo userRepo.Repository
}

func NewService(repo userRepo.Repository) Service {
	return &service{repo: repo}
}

func (s *service) Register(ctx context.Context, user *userModels.User) (*userModels.User, error) {
	if user.Role == "superadmin" {
		return nil, fmt.Errorf("superadmin cannot be created")
	}
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) GetProfile(ctx context.Context, id uint) (*userModels.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// TODO: Implement the UpdateProfile method (Remove password field from the user struct)
func (s *service) UpdateProfile(ctx context.Context, user *userModels.User) error {
	if user.ID == 0 {
		return errors.New("user ID is required")
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteUser(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("user ID is required")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

func (s *service) ValidateUser(ctx context.Context, email, password string) (bool, *userModels.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return false, nil, err
	}

	valid, err := user.ComparePassword(password)

	if err != nil {
		return false, nil, err
	}

	if valid {
		return true, user, nil
	}

	return false, nil, nil
}
