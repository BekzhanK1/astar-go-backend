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

func (s *AuthService) Login(email, password string) (string, string, string, error) {
	fmt.Printf("Validating user %s\n", email)
	fmt.Printf("Validating user %s\n", password)
	userID, role, err := s.userClient.ValidateUser(email, password)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to validate user: %w", err)
	}

	accessToken, refreshToken, err := s.service.GenerateTokens(uint(userID), role)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate tokens: %w", err)
	}

	hashToken, err := utils.HashToken(refreshToken)

	if err != nil {
		return "", "", "", err
	}

	sessionID, err := s.SaveRefreshToken(uint(userID), hashToken)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to save refresh token: %w", err)
	}

	return accessToken, refreshToken, sessionID, nil
}

func (s *AuthService) SaveRefreshToken(userID uint, tokenHash string) (string, error) {
	token := models.RefreshToken{
		UserID:    userID,
		TokenHash: tokenHash,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	log.Println("Saving refresh token to Redis")
	return s.repo.SaveRefreshToken(userID, token, 7*24*time.Hour)
}

func (s *AuthService) ValidateToken(token string) (uint64, string, bool, error) {
	claims, err := s.service.ValidateToken(token)
	if err != nil {
		return 0, "", false, fmt.Errorf("invalid token: %w", err)
	}

	userID := uint64(claims["user_id"].(float64))
	role := claims["role"].(string)
	return userID, role, true, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	claims, err := s.service.ValidateToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	userID := uint(claims["user_id"].(float64))
	role := claims["role"].(string)

	// Generate a new access token
	accessToken, _, err := s.service.GenerateTokens(userID, role)
	if err != nil {
		return "", fmt.Errorf("failed to generate new access token: %w", err)
	}

	return accessToken, nil
}

func (s *AuthService) Logout(userID uint, sessionID string) error {
	if err := s.repo.DeleteRefreshToken(userID, sessionID); err != nil {
		return fmt.Errorf("failed to delete refresh token for session %s: %w", sessionID, err)
	}
	return nil
}

func (s *AuthService) RevokeAllTokens(userID uint) error {
	if err := s.repo.RevokeAllTokens(userID); err != nil {
		return fmt.Errorf("failed to revoke all tokens for user %d: %w", userID, err)
	}
	return nil
}
