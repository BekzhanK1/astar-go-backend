package main

import (
	authGrpc "auth-service/api/grpc"
	"auth-service/clients"
	"auth-service/internal/repositories"
	"auth-service/internal/services"
	"auth-service/internal/utils"
	"auth-service/pkg/auth"
	"auth-service/pkg/config"
	"log"
	"net"

	pb "auth-service/protobuf/auth"

	"google.golang.org/grpc"
)

func main() {
	utils.LoadEnv()
	tokenService := auth.NewAuthService(utils.GetEnv("JWT_SECRET", "your-secret-key"))
	userClient, err := clients.NewUserServiceClient(utils.GetEnv("USER_SERVICE_URL", "localhost:50051"))
	if err != nil {
		log.Fatalf("Failed to create user client: %v", err)
	}

	redisClient := config.NewConfig().RedisClient
	tokenRepo := repositories.NewTokenRepository(redisClient) // Replace with actual DB instance
	authService := services.NewAuthService(tokenService, userClient, tokenRepo)
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
