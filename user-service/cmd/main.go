package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	userGrpc "user-service/api/grpc"
	"user-service/internal/user"
	"user-service/pkg/config"
	"user-service/pkg/database"
	pb "user-service/proto"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.InitDB(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	// Auto-migrate User schema
	if err := db.AutoMigrate(&user.User{}); err != nil {
		log.Fatalf("failed to auto-migrate: %v", err)
	}

	// Initialize repository and service
	userRepo := user.NewGormRepository(db)
	userService := user.NewService(userRepo)

	// Set up gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userGrpc.NewUserServiceServer(userService))

	// Listen and serve
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen on port 50051: %v", err)
	}

	log.Println("User service is running on port :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC server: %v", err)
	}
}
