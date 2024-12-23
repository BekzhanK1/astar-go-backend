package grpc

import (
	"auth-service/internal/services"
	pb "auth-service/protobuf/auth"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	accessToken, refreshToken, sessionID, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to login: %v", err)
	}

	return &pb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		SessionId:    sessionID,
	}, nil
}

func (h *AuthHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	userID, role, valid, err := h.authService.ValidateToken(req.AccessToken)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to validate token: %v", err)
	}

	return &pb.ValidateTokenResponse{
		UserId: userID,
		Role:   role,
		Valid:  valid,
	}, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	accessToken, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to refresh token: %v", err)
	}

	return &pb.RefreshTokenResponse{
		AccessToken: accessToken,
	}, nil
}

func (h *AuthHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	userID, sessionID := uint(req.UserId), req.SessionId

	if userID == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID")
	}

	if sessionID == "" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid session ID")
	}

	err := h.authService.Logout(userID, sessionID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to logout: %v", err)
	}

	return &pb.LogoutResponse{
		Message: "Logout successful",
	}, nil
}

func (h *AuthHandler) RevokeAllTokens(ctx context.Context, req *pb.RevokeAllTokensRequest) (*pb.RevokeAllTokensResponse, error) {
	err := h.authService.RevokeAllTokens(uint(req.UserId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to revoke all tokens: %v", err)
	}

	return &pb.RevokeAllTokensResponse{
		Message: "All tokens revoked successfully",
	}, nil
}
