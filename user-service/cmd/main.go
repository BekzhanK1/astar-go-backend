package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	userGrpc "user-service/api/grpc"
	admin "user-service/internal/admin"
	userModel "user-service/internal/user/models"
	userRepository "user-service/internal/user/repository"
	userService "user-service/internal/user/service"
	"user-service/pkg/config"
	"user-service/pkg/database"
	pb "user-service/protobuf/user"
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
	if err = db.AutoMigrate(&userModel.User{}); err != nil {
		log.Fatalf("failed to auto-migrate: %v", err)
	}

	// Initialize repository and service
	userRepo := userRepository.NewGormRepository(db)
	userService := userService.NewService(userRepo)

	adminUser, err := admin.CreateAdmin(userRepo)

	if err != nil {
		log.Fatalf("failed to create admin user: %v", err)
	}

	log.Printf("Admin user created: %s", adminUser.Email)

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
