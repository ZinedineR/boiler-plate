package handler

import (
	"boiler-plate/internal/transaction/service"
	transaction "boiler-plate/proto/transaction/v1"
	"boiler-plate/proto/users/v1"
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type GRPCHandler struct {
	transaction.UnimplementedServiceServer
	TransactionService service.Service
}

func NewGRPCHandler(service service.Service) *GRPCHandler {
	return &GRPCHandler{TransactionService: service}
}

func (s *GRPCHandler) CreditTransaction(
	ctx context.Context, in *transaction.CreditTransactionRequest,
) (*transaction.CreditTransactionResponse, error) {
	if err := s.TransactionService.Credit(ctx, in.GetUserid(), in.GetAmount()); err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &transaction.CreditTransactionResponse{
		Response: &transaction.MutationResponse{Message: "Credit to user: " + in.GetUserid() + "successfully"},
	}, nil
}

func (s *GRPCHandler) GetTransaction(
	ctx context.Context, in *transaction.GetTransactionRequest,
) (*transaction.GetTransactionResponse, error) {
	if in.GetLimit() == "" {
		in.Limit = "0"
	}
	if in.GetPage() == "" {
		in.Page = "0"
	}
	result, err := s.TransactionService.Find(ctx, in.GetLimit(), in.GetPage(), in.GetUserid())
	if err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	var transactionProto []*transaction.Transactions
	for _, dataTransaction := range result.Data {
		transactionProto = append(transactionProto, &transaction.Transactions{
			Id:              int32(dataTransaction.ID),
			Type:            dataTransaction.Type,
			Amount:          dataTransaction.Amount,
			Message:         dataTransaction.Message,
			TransactionTime: timestamppb.New(*dataTransaction.TransactionTime),
		})
	}
	return &transaction.GetTransactionResponse{
		Pagination: &transaction.PaginationResponse{
			Limit:      int32(result.Pagination.Limit),
			Page:       int32(result.Pagination.Page),
			TotalRows:  int32(result.Pagination.TotalRows),
			TotalPages: int32(result.Pagination.TotalPages),
		},
		Users: &users.Users{
			Id:        int32(result.Users.ID),
			Name:      result.Users.Name,
			Email:     result.Users.Email,
			Password:  result.Users.Password,
			CreatedAt: timestamppb.New(*result.Users.CreatedAt),
			UpdatedAt: timestamppb.New(*result.Users.UpdatedAt),
			WalletId:  int32(result.Users.WalletId),
			Wallet: &users.Wallet{
				Id:              int32(result.Users.Wallet.ID),
				Balance:         result.Users.Wallet.Balance,
				LastTransaction: timestamppb.New(*result.Users.Wallet.LastTransaction),
			},
		},
		Transaction: transactionProto,
		Response:    &transaction.MutationResponse{Message: "Find Transaction Success"},
	}, nil
}

func (s *GRPCHandler) TransferTransaction(
	ctx context.Context, in *transaction.TransferTransactionRequest,
) (*transaction.TransferTransactionResponse, error) {
	if err := s.TransactionService.Transfer(ctx, in.GetSenderid(), in.GetReceiverid(), in.GetAmount()); err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &transaction.TransferTransactionResponse{
		Response: &transaction.MutationResponse{Message: "Transfer Success from id: " + in.GetSenderid() + "to id: " + in.GetReceiverid() + "successfully"},
	}, nil
}

func (s *GRPCHandler) DetailTransaction(
	ctx context.Context, in *transaction.DetailTransactionRequest,
) (*transaction.DetailTransactionResponse, error) {
	result, err := s.TransactionService.Detail(ctx, in.GetId())
	if err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &transaction.DetailTransactionResponse{
		Transaction: &transaction.Transactions{
			Id:      int32(result.ID),
			Type:    result.Type,
			Amount:  result.Amount,
			Message: result.Message,
			Userid:  wrapperspb.Int32(int32(result.Users.ID)),
			Users: &users.Users{
				Id:        int32(result.Users.ID),
				Name:      result.Users.Name,
				Email:     result.Users.Email,
				Password:  result.Users.Password,
				CreatedAt: timestamppb.New(*result.Users.CreatedAt),
				UpdatedAt: timestamppb.New(*result.Users.UpdatedAt),
			},
			TransactionTime: timestamppb.New(*result.TransactionTime),
		},
		Response: &transaction.MutationResponse{Message: "Find Transaction Success"},
	}, nil
}
