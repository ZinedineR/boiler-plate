package service

import (
	"boiler-plate/internal/users/domain"
	"boiler-plate/pkg/db"
	"boiler-plate/pkg/exception"
	"context"
)

type Service interface {
	Create(
		ctx context.Context, req *domain.Users,
	) *exception.Exception
	// Detail Service
	Detail(ctx context.Context, id string) (*domain.UserResponse, *exception.Exception)
	// Delete Service
	Delete(ctx context.Context, id string) *exception.Exception
	// Update Service
	Update(
		ctx context.Context, id string, users *domain.Users,
	) *exception.Exception
	Find(ctx context.Context, limit, page string) (*FindResponse, *exception.Exception)
	Auth(ctx context.Context, email, password string) (*domain.Users, *exception.Exception)
}

type FindResponse struct {
	Pagination db.Paginate           `json:"pagination"`
	Data       []domain.UserResponse `json:"data"`
}
