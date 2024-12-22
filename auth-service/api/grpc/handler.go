package grpc

import (
	"auth-service/internal/services"
	pb "auth-service/proto"
	"context"
	"fmt"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	accessToken, refreshToken, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}
	return &pb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
