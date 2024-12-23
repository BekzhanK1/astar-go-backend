package clients

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "auth-service/protobuf/auth"

	"google.golang.org/grpc"
)

type UserServiceClient struct {
	client pb.UserServiceClient
}

func NewUserServiceClient(address string) (*UserServiceClient, error) {
	log.Printf("Connecting to user-service at %s", address)

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user-service: %w", err)
	}

	log.Println("Successfully connected to user-service")
	return &UserServiceClient{
		client: pb.NewUserServiceClient(conn),
	}, nil
}

func (a *UserServiceClient) ValidateUser(email, password string) (uint64, string, error) {
	fmt.Printf("ValidateUser: Initiating request for email: %s\n", email)

	req := &pb.ValidateUserRequest{
		Email:    email,
		Password: password,
	}

	fmt.Println("ValidateUser: Making gRPC call to user-service")

	resp, err := a.client.ValidateUser(context.Background(), req)
	if err != nil {
		fmt.Printf("ValidateUser: gRPC call failed with error: %v\n", err)
		return 0, "", fmt.Errorf("ValidateUser: error making gRPC call: %w", err)
	}

	fmt.Println("ValidateUser: gRPC call successful")

	if !resp.Valid {
		fmt.Printf("ValidateUser: Invalid credentials for email: %s\n", email)
		return 0, "", fmt.Errorf("ValidateUser: invalid user credentials for email: %s", email)
	}

	fmt.Printf("ValidateUser: Successfully validated user ID %d with role %s\n", resp.Id, resp.Role)
	return resp.Id, resp.Role, nil
}
