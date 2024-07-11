package service

import (
	"boiler-plate/internal/users/domain"
	"boiler-plate/pkg/db"
	"boiler-plate/pkg/exception"
	"context"
)

type Service interface {
	Create(
		ctx context.Context, req *UserRequest,
	) *exception.Exception
	// Detail Service
	Detail(ctx context.Context, id string) (*domain.Users, *exception.Exception)
	// Delete Service
	Delete(ctx context.Context, id string) *exception.Exception
	// Update Service
	Update(
		ctx context.Context, id string, users *UserRequest,
	) *exception.Exception
	Find(ctx context.Context, limit, page string) (*FindResponse, *exception.Exception)
}

type UserRequest struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `validate:"required,gt=2" json:"email,omitempty"`
	Password string `validate:"required" json:"password,omitempty"`
}

type FindResponse struct {
	Pagination db.Paginate    `json:"pagination"`
	Data       []domain.Users `json:"data"`
}
