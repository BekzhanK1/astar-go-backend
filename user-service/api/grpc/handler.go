package grpc

import (
	"context"

	"user-service/internal/user"
	pb "user-service/proto"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	service user.Service
}

// NewUserServiceServer creates a new gRPC server for the UserService.
func NewUserServiceServer(service user.Service) *UserServiceServer {
	return &UserServiceServer{service: service}
}

func (s *UserServiceServer) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	user := &user.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
		Password:  req.Password,
	}

	createdUser, err := s.service.Register(ctx, user)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterUserResponse{
		User: &pb.User{
			Id:        uint64(createdUser.ID),
			Email:     createdUser.Email,
			FirstName: createdUser.FirstName,
			LastName:  createdUser.LastName,
			Role:      createdUser.Role,
			CreatedAt: createdUser.CreatedAt,
			UpdatedAt: createdUser.UpdatedAt,
		},
	}, nil
}

func (s *UserServiceServer) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	userProfile, err := s.service.GetProfile(ctx, uint(req.Id))
	if err != nil {
		return nil, err
	}

	return &pb.GetProfileResponse{
		User: &pb.User{
			Id:        uint64(userProfile.ID),
			Email:     userProfile.Email,
			FirstName: userProfile.FirstName,
			LastName:  userProfile.LastName,
			Role:      userProfile.Role,
			CreatedAt: userProfile.CreatedAt,
			UpdatedAt: userProfile.UpdatedAt,
		},
	}, nil
}

func (s *UserServiceServer) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	user := &user.User{
		ID:        uint(req.User.Id),
		Email:     req.User.Email,
		FirstName: req.User.FirstName,
		LastName:  req.User.LastName,
		Role:      req.User.Role,
		Password:  req.User.Password,
	}

	if err := s.service.UpdateProfile(ctx, user); err != nil {
		return nil, err
	}

	return &pb.UpdateProfileResponse{
		User: req.User,
	}, nil
}

func (s *UserServiceServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if err := s.service.DeleteUser(ctx, uint(req.Id)); err != nil {
		return nil, err
	}

	return &pb.DeleteUserResponse{
		Success: true,
	}, nil
}

func (s *UserServiceServer) ValidateUser(ctx context.Context, req *pb.ValidateUserRequest) (*pb.ValidateUserResponse, error) {
	valid, err := s.service.ValidateUser(ctx, req.Email, req.Password)

	if err != nil {
		return nil, err
	}

	return &pb.ValidateUserResponse{
		Valid: valid,
	}, nil
}
