package services

import (
	"auth-service/clients"
	"auth-service/internal/models"
	"auth-service/internal/repositories"
	"auth-service/internal/utils"
	"auth-service/pkg/auth"
	"fmt"
	"log"
	"time"
)

type UserRole string

const (
	SuperAdmin  UserRole = "superadmin"
	SchoolAdmin UserRole = "school_admin"
	Advisor     UserRole = "advisor"
	Teacher     UserRole = "teacher"
)

type AuthService struct {
	service    *auth.TokenService
	userClient *clients.UserServiceClient
	repo       *repositories.TokenRepository
}

func NewAuthService(service *auth.TokenService, userClient *clients.UserServiceClient, repo *repositories.TokenRepository) *AuthService {
	return &AuthService{service: service, userClient: userClient, repo: repo}
}

func (s *AuthService) Login(email, password string) (string, string, error) {
	fmt.Printf("Validating user %s\n", email)
	fmt.Printf("Validating user %s\n", password)
	userID, role, err := s.userClient.ValidateUser(email, password)
	if err != nil {
		return "", "", fmt.Errorf("failed to validate user: %w", err)
	}

	accessToken, refreshToken, err := s.service.GenerateTokens(uint(userID), role)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate tokens: %w", err)
	}

	hashToken, err := utils.HashToken(refreshToken)

	if err != nil {
		return "", "", err
	}

	if err := s.SaveRefreshToken(uint(userID), hashToken); err != nil {
		return "", "", fmt.Errorf("failed to save refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) SaveRefreshToken(userID uint, tokenHash string) error {
	token := models.RefreshToken{
		UserID:    userID,
		TokenHash: tokenHash,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	log.Println("Saving refresh token to Redis")
	return s.repo.SaveRefreshToken(userID, token, 7*24*time.Hour)
}

// func (s *AuthService) ValidateRefreshToken(userID uint, tokenHash string) (bool, error) {
// 	valid, err := s.repo.ValidateRefreshToken(userID, tokenHash)
// 	if err != nil {
// 		return false, fmt.Errorf("token validation failed: %w", err)
// 	}
// 	return valid, nil
// }

// func (s *AuthService) RevokeRefreshToken(userID uint) error {
// 	return s.repo.DeleteRefreshToken(userID)
// }
