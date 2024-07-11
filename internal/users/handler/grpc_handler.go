package handler

import (
	"boiler-plate/internal/users/domain"
	"boiler-plate/internal/users/service"
	users "boiler-plate/proto/users/v1"
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCHandler struct {
	users.UnimplementedServiceServer
	UsersService service.Service
}

func NewGRPCHandler(service service.Service) *GRPCHandler {
	return &GRPCHandler{UsersService: service}
}

func (s *GRPCHandler) CreateUser(ctx context.Context, in *users.CreateUserRequest) (*users.CreateUserResponse, error) {
	body := &domain.Users{
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	}
	if err := s.UsersService.Create(ctx, body); err != nil {
		return nil, err.Error
	}
	return &users.CreateUserResponse{
		Data: &users.Users{
			Id:        int32(body.ID),
			Email:     in.Email,
			Password:  in.Password,
			CreatedAt: timestamppb.New(*body.CreatedAt),
			UpdatedAt: timestamppb.New(*body.UpdatedAt),
		},
		Response: &users.MutationResponse{Message: "Create User Success"},
	}, nil
}

func (s *GRPCHandler) GetUser(ctx context.Context, in *users.GetUserRequest) (*users.GetUserResponse, error) {

	result, err := s.UsersService.Find(ctx, in.GetLimit(), in.GetPage())
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("get user error: %v", err))
	}
	var usersProto []*users.Users
	for _, dataUser := range result.Data {
		usersProto = append(usersProto, &users.Users{
			Id:        int32(dataUser.ID),
			Email:     dataUser.Email,
			Password:  dataUser.Password,
			CreatedAt: timestamppb.New(*dataUser.CreatedAt),
			UpdatedAt: timestamppb.New(*dataUser.UpdatedAt),
		})
	}
	return &users.GetUserResponse{
		Pagination: &users.PaginationResponse{
			Limit:      int32(result.Pagination.Limit),
			Page:       int32(result.Pagination.Page),
			TotalRows:  int32(result.Pagination.TotalRows),
			TotalPages: int32(result.Pagination.TotalPages),
		},
		Users:    usersProto,
		Response: &users.MutationResponse{Message: "Find User Success"},
	}, nil
}

func (s *GRPCHandler) UpdateUser(ctx context.Context, in *users.UpdateUserRequest) (*users.UpdateUserResponse, error) {
	body := &domain.Users{
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	}
	if err := s.UsersService.Update(ctx, in.GetId(), body); err != nil {
		return nil, err.Error
	}
	return &users.UpdateUserResponse{
		Data: &users.Users{
			Id:        int32(body.ID),
			Email:     in.Email,
			Password:  in.Password,
			CreatedAt: timestamppb.New(*body.CreatedAt),
			UpdatedAt: timestamppb.New(*body.UpdatedAt),
		},
		Response: &users.MutationResponse{Message: "Create User Success"},
	}, nil
}

func (s *GRPCHandler) DetailUser(ctx context.Context, in *users.DetailUserRequest) (*users.DetailUserResponse, error) {
	result, err := s.UsersService.Detail(ctx, in.GetId())
	if err != nil {
		return nil, err.Error
	}
	return &users.DetailUserResponse{
		User: &users.Users{
			Id:        int32(result.ID),
			Email:     result.Email,
			Password:  result.Password,
			CreatedAt: timestamppb.New(*result.CreatedAt),
			UpdatedAt: timestamppb.New(*result.UpdatedAt),
		},
	}, nil
}

func (s *GRPCHandler) DeleteUser(ctx context.Context, in *users.DeleteUserRequest) (*users.DeleteUserResponse, error) {
	if err := s.UsersService.Delete(ctx, in.GetId()); err != nil {
		return nil, err.Error
	}
	return &users.DeleteUserResponse{
		Response: &users.MutationResponse{Message: "Delete User Success"},
	}, nil
}
