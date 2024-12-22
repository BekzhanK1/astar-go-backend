package services

import (
	"auth-service/internal/models"
	"auth-service/internal/repositories"
	"auth-service/internal/utils"
	"auth-service/pkg/auth"
	"fmt"
	"time"
)

type AuthService struct {
	service *auth.TokenService
	repo    *repositories.TokenRepository
}

func NewAuthService(repo *repositories.TokenRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Login(userID uint, role ) (string, string, error) {
	accessToken, refreshToken, err := s.service.GenerateTokens(userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate tokens: %w", err)
	}

	tokenHash, err := utils.HashToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to hash token: %w", err)
	}

	if err := s.SaveRefreshToken(userID, tokenHash); err != nil {
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
	return s.repo.SaveRefreshToken(userID, token, 7*24*time.Hour)
}

func (s *AuthService) ValidateRefreshToken(userID uint, tokenHash string) (bool, error) {
	valid, err := s.repo.ValidateRefreshToken(userID, tokenHash)
	if err != nil {
		return false, fmt.Errorf("token validation failed: %w", err)
	}
	return valid, nil
}

func (s *AuthService) RevokeRefreshToken(userID uint) error {
	return s.repo.DeleteRefreshToken(userID)
}
