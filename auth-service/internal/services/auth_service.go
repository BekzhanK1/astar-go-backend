package services

import (
	"auth-service/clients"
	"auth-service/pkg/auth"
	"fmt"
)

type UserRole string

const (
	SuperAdmin  UserRole = "superadmin"
	SchoolAdmin UserRole = "school_admin"
	Advisor     UserRole = "advisor"
	Teacher     UserRole = "teacher"
)

type AuthService struct {
	service *auth.TokenService
	// repo       *repositories.TokenRepository
	userClient *clients.UserServiceClient
}

func NewAuthService(service *auth.TokenService, userClient *clients.UserServiceClient) *AuthService {
	return &AuthService{service: service, userClient: userClient}
}

func (s *AuthService) Login(email, password string) (string, string, error) {
	// Step 1: Validate user credentials via gRPC call
	fmt.Printf("Validating user %s\n", email)
	fmt.Printf("Validating user %s\n", password)
	userID, role, err := s.userClient.ValidateUser(email, password)
	if err != nil {
		return "", "", fmt.Errorf("failed to validate user: %w", err)
	}

	// Step 2: Generate access and refresh tokens
	accessToken, refreshToken, err := s.service.GenerateTokens(uint(userID), role)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate tokens: %w", err)
	}

	return accessToken, refreshToken, nil
}

// func (s *AuthService) SaveRefreshToken(userID uint, tokenHash string) error {
// 	token := models.RefreshToken{
// 		UserID:    userID,
// 		TokenHash: tokenHash,
// 		IssuedAt:  time.Now(),
// 		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
// 	}
// 	return s.repo.SaveRefreshToken(userID, token, 7*24*time.Hour)
// }

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
