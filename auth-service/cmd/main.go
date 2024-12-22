package main

import (
	authGrpc "auth-service/api/grpc"
	"auth-service/clients"
	"auth-service/internal/services"
	"auth-service/pkg/auth"
	"log"
	"net"

	pb "auth-service/proto"

	"google.golang.org/grpc"
)

func main() {
	// Initialize dependencies
	tokenService := auth.NewAuthService("your-secret-key")
	// tokenRepo := repositories.NewTokenRepository(nil)                  // Replace with actual DB instance
	userClient, err := clients.NewUserServiceClient("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to create user client: %v", err)
	}

	authService := services.NewAuthService(tokenService, userClient)
	authHandler := authGrpc.NewAuthHandler(authService)

	// Set up gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authHandler)

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen on port 50052: %v", err)
	}

	log.Println("Auth service running on port 50052")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
