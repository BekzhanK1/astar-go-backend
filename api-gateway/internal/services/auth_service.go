package services

import (
	pb "api-gateway/protobuf/auth"
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
)

type AuthServiceClient struct {
	client pb.AuthServiceClient
}

func NewAuthServiceClient(address string) (*AuthServiceClient, error) {
	log.Printf("Connecting to auth-service at %s", address)

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth-service: %w", err)
	}

	log.Println("Successfully connected to auth-service")
	return &AuthServiceClient{
		client: pb.NewAuthServiceClient(conn),
	}, nil
}

func (a *AuthServiceClient) Login(email, password string) (string, string, string, error) {
	fmt.Printf("Login: Initiating request for email: %s\n", email)

	req := &pb.LoginRequest{
		Email:    email,
		Password: password,
	}

	fmt.Println("Login: Making gRPC call to auth-service")

	resp, err := a.client.Login(context.Background(), req)
	if err != nil {
		fmt.Printf("Login: gRPC call failed with error: %v\n", err)
		return "", "", "", fmt.Errorf("Login: error making gRPC call: %w", err)
	}

	fmt.Println("Login: gRPC call successful")

	fmt.Printf("Login: Successfully logged in with session ID %s\n", resp.SessionId)
	return resp.AccessToken, resp.RefreshToken, resp.SessionId, nil
}

func (a *AuthServiceClient) ValidateToken(accessToken string) (uint64, string, bool, error) {
	fmt.Printf("ValidateToken: Initiating request for access token: %s\n", accessToken)

	req := &pb.ValidateTokenRequest{
		AccessToken: accessToken,
	}

	fmt.Println("ValidateToken: Making gRPC call to auth-service")

	resp, err := a.client.ValidateToken(context.Background(), req)
	if err != nil {
		fmt.Printf("ValidateToken: gRPC call failed with error: %v\n", err)
		return 0, "", false, fmt.Errorf("ValidateToken: error making gRPC call: %w", err)
	}

	fmt.Println("ValidateToken: gRPC call successful")

	fmt.Printf("ValidateToken: Successfully validated token for user ID %d with role %s\n", resp.UserId, resp.Role)
	return resp.UserId, resp.Role, resp.Valid, nil
}
