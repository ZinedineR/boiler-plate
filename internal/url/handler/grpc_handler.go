package handler

import (
	"boiler-plate/internal/url/domain"
	"boiler-plate/internal/url/service"
	"boiler-plate/proto/url/v1"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCHandler struct {
	url.UnimplementedServiceServer
	URLService service.Service
}

func NewGRPCHandler(service service.Service) *GRPCHandler {
	return &GRPCHandler{URLService: service}
}

func (s *GRPCHandler) CreateURL(ctx context.Context, in *url.CreateURLRequest) (*url.CreateURLResponse, error) {
	body := &domain.URL{
		Longurl: in.GetLongurl(),
	}
	if err := s.URLService.Create(ctx, body); err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &url.CreateURLResponse{
		Data: &url.URLs{
			Id:        int32(body.ID),
			Longurl:   body.Longurl,
			Shorturl:  body.Shorturl,
			CreatedAt: timestamppb.New(*body.CreatedAt),
		},
		Response: &url.MutationResponse{Message: "Create User Success"},
	}, nil
}

func (s *GRPCHandler) DetailURL(ctx context.Context, in *url.DetailURLRequest) (*url.DetailURLResponse, error) {

	result, err := s.URLService.Detail(ctx, in.GetShorturl())
	if err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	header := metadata.Pairs("Location", result.Longurl)
	errGrpc := grpc.SendHeader(ctx, header)
	if errGrpc != nil {
		return nil, status.Error(codes.Code(16), fmt.Sprint("error sending header"))
	}
	return &url.DetailURLResponse{
		Url: &url.URLs{
			Id:        int32(result.ID),
			Longurl:   result.Longurl,
			Shorturl:  result.Shorturl,
			CreatedAt: timestamppb.New(*result.CreatedAt),
		},
	}, nil
}

//
//func (s *GRPCHandler) DeleteUser(ctx context.Context, in *url.DeleteUserRequest) (*url.DeleteUserResponse, error) {
//	if err := s.UsersService.Delete(ctx, in.GetId()); err != nil {
//		return nil, err.Error
//	}
//	return &url.DeleteUserResponse{
//		Response: &url.MutationResponse{Message: "Delete User Success"},
//	}, nil
//}
