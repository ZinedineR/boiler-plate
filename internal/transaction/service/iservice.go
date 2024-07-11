package service

import (
	"boiler-plate/internal/transaction/domain"
	userDomain "boiler-plate/internal/users/domain"
	"boiler-plate/pkg/db"
	"boiler-plate/pkg/exception"
	"context"
)

type Service interface {
	Credit(
		ctx context.Context, userid string, amount float64,
	) *exception.Exception
	Transfer(
		ctx context.Context, senderId, receiverId string, amount float64,
	) *exception.Exception
	Detail(ctx context.Context, id string) (*domain.Transaction, *exception.Exception)
	Find(ctx context.Context, limit, page, userid string) (*FindResponse, *exception.Exception)
}

type FindResponse struct {
	Pagination db.Paginate          `json:"pagination"`
	Users      userDomain.Users     `json:"users"`
	Data       []domain.Transaction `json:"data"`
}
