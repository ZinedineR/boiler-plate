package handler

import (
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
	body := &service.UserRequest{
		Name:     in.GetName(),
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	}
	if err := s.UsersService.Create(ctx, body); err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &users.CreateUserResponse{
		Data: &users.Users{
			Id:       int32(body.Id),
			Name:     in.Name,
			Email:    in.Email,
			Password: in.Password,
		},
		Response: &users.MutationResponse{Message: "Create User Success"},
	}, nil
}

func (s *GRPCHandler) GetUser(ctx context.Context, in *users.GetUserRequest) (*users.GetUserResponse, error) {
	if in.GetLimit() == "" {
		in.Limit = "0"
	}
	if in.GetPage() == "" {
		in.Page = "0"
	}
	result, err := s.UsersService.Find(ctx, in.GetLimit(), in.GetPage())
	if err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	var usersProto []*users.Users
	for _, dataUser := range result.Data {
		usersProto = append(usersProto, &users.Users{
			Id:       int32(dataUser.ID),
			Name:     dataUser.Name,
			Email:    dataUser.Email,
			Password: dataUser.Password,
			WalletId: int32(dataUser.WalletId),
			Wallet: &users.Wallet{
				Id:              int32(dataUser.Wallet.ID),
				Balance:         dataUser.Wallet.Balance,
				LastTransaction: timestamppb.New(*dataUser.Wallet.LastTransaction),
			},
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
	body := &service.UserRequest{
		Name:     in.GetName(),
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	}
	if err := s.UsersService.Update(ctx, in.GetId(), body); err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &users.UpdateUserResponse{
		Data: &users.Users{
			Id:       int32(body.Id),
			Name:     in.Name,
			Email:    in.Email,
			Password: in.Password,
		},
		Response: &users.MutationResponse{Message: "Create User Success"},
	}, nil
}

func (s *GRPCHandler) DetailUser(ctx context.Context, in *users.DetailUserRequest) (*users.DetailUserResponse, error) {
	result, err := s.UsersService.Detail(ctx, in.GetId())
	if err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &users.DetailUserResponse{
		User: &users.Users{
			Id:       int32(result.ID),
			Name:     result.Name,
			Email:    result.Email,
			Password: result.Password,
			WalletId: int32(result.WalletId),
			Wallet: &users.Wallet{
				Id:              int32(result.Wallet.ID),
				Balance:         result.Wallet.Balance,
				LastTransaction: timestamppb.New(*result.Wallet.LastTransaction),
			},
			CreatedAt: timestamppb.New(*result.CreatedAt),
			UpdatedAt: timestamppb.New(*result.UpdatedAt),
		},
	}, nil
}

func (s *GRPCHandler) DeleteUser(ctx context.Context, in *users.DeleteUserRequest) (*users.DeleteUserResponse, error) {
	if err := s.UsersService.Delete(ctx, in.GetId()); err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &users.DeleteUserResponse{
		Response: &users.MutationResponse{Message: "Delete User Success"},
	}, nil
}
