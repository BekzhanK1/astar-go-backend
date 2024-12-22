package clients

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "auth-service/proto" // Correctly imports the auth-service proto package

	"google.golang.org/grpc"
)

type UserServiceClient struct { // Updated struct name to reflect the service
	client pb.UserServiceClient
}

func NewUserServiceClient(address string) (*UserServiceClient, error) {
	log.Printf("Connecting to user-service at %s", address) // Debug log

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user-service: %w", err)
	}

	log.Println("Successfully connected to user-service") // Debug log
	return &UserServiceClient{
		client: pb.NewUserServiceClient(conn),
	}, nil
}

func (a *UserServiceClient) ValidateUser(email, password string) (uint64, string, error) {
	fmt.Printf("ValidateUser: Initiating request for email: %s\n", email)

	// Build the request
	req := &pb.ValidateUserRequest{
		Email:    email,
		Password: password,
	}

	fmt.Println("ValidateUser: Making gRPC call to user-service") // Debug log

	// Make the gRPC call without a timeout
	resp, err := a.client.ValidateUser(context.Background(), req)
	if err != nil {
		fmt.Printf("ValidateUser: gRPC call failed with error: %v\n", err) // Debug log
		return 0, "", fmt.Errorf("ValidateUser: error making gRPC call: %w", err)
	}

	fmt.Println("ValidateUser: gRPC call successful") // Debug log

	// Check if the user is valid
	if !resp.Valid {
		fmt.Printf("ValidateUser: Invalid credentials for email: %s\n", email) // Debug log
		return 0, "", fmt.Errorf("ValidateUser: invalid user credentials for email: %s", email)
	}

	fmt.Printf("ValidateUser: Successfully validated user ID %d with role %s\n", resp.Id, resp.Role)
	// Return user ID and role
	return resp.Id, resp.Role, nil
}
