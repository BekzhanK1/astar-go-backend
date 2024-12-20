package grpc

import (
	pb "user-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewGRPCServer(userService pb.UserServiceServer) *grpc.Server {
	server := grpc.NewServer()

	// Register the user service
	pb.RegisterUserServiceServer(server, userService)

	// Enable reflection for debugging with tools like grpcurl
	reflection.Register(server)

	return server
}
