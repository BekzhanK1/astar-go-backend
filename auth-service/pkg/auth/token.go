package auth

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

type TokenService struct {
	JWTSecret []byte
}

func NewAuthService(secret string) *TokenService {
	return &TokenService{
		JWTSecret: []byte(secret),
	}
}

// GenerateTokens generates an access token and a refresh token with user ID and role.
func (ts *TokenService) GenerateTokens(userID uint, role string) (string, string, error) {
	// Access Token
	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(ts.JWTSecret)
	if err != nil {
		return "", "", err
	}

	// Refresh Token
	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(ts.JWTSecret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ValidateToken validates a JWT token and returns the claims.
func (ts *TokenService) ValidateToken(token string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return ts.JWTSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}
