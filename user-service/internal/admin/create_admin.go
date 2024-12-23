package admin

import (
	"context"
	userModels "user-service/internal/user/models"
	userRepo "user-service/internal/user/repository"
	"user-service/internal/utils"
)

func CreateAdmin(userRepo userRepo.Repository) (*userModels.User, error) {
	hashedPassword, err := utils.HashPassword(utils.GetEnv("ADMIN_PASSWORD", "superadmin"))

	if err != nil {
		return nil, err
	}

	user := &userModels.User{
		Email:     utils.GetEnv("ADMIN_EMAIL", "superadmin@gmail.com"),
		FirstName: "Admin",
		LastName:  "Admin",
		Role:      "superadmin",
		Password:  hashedPassword,
	}

	ctx := context.Background()

	existingUser, err := userRepo.GetByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return existingUser, nil
	}

	if err := userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
