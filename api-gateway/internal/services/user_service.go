package services

import (
	pb "api-gateway/protobuf/user"
	"context"
	"fmt"
	"log"
	"time"

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

func (u *UserServiceClient) GetProfile(userID uint64) (*pb.GetProfileResponse, error) {
	req := &pb.GetProfileRequest{
		Id: userID,
	}

	resp, err := u.client.GetProfile(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	return resp, nil
}
