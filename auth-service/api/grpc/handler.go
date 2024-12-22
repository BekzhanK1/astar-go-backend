package handlers

import (
	"auth-service/internal/services"
	pb "auth-service/proto"
	"context"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authService services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return h.authService.Login(req.Email, req.Password)
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	return h.authService.RefreshToken(req.RefreshToken)
}
