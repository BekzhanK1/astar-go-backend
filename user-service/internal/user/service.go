package user

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, user *User) (*User, error)
	GetProfile(ctx context.Context, id uint) (*User, error)
	UpdateProfile(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Register(ctx context.Context, user *User) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) GetProfile(ctx context.Context, id uint) (*User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) UpdateProfile(ctx context.Context, user *User) error {
	if user.ID == 0 {
		return errors.New("user ID is required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
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
