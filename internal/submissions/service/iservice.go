package service

import (
	"boiler-plate/internal/submissions/domain"
	"boiler-plate/pkg/db"
	"boiler-plate/pkg/exception"
	"context"
)

type Service interface {
	Create(
		ctx context.Context, req *domain.SubmissionRequest,
	) *exception.Exception
	// Detail Service
	Detail(ctx context.Context, id string) (*domain.Submissions, *exception.Exception)
	// DetailByUser Service
	Delete(ctx context.Context, id string) *exception.Exception
	Find(ctx context.Context, limit, page string) (*FindResponse, *exception.Exception)
	FindByUser(ctx context.Context, limit, page, userid string) (*FindByUserResponse, *exception.Exception)
}

type FindByUserResponse struct {
	UserId     int                  `json:"user_id"`
	Pagination db.Paginate          `json:"pagination"`
	Data       []domain.Submissions `json:"data"`
}

type FindResponse struct {
	Pagination db.Paginate          `json:"pagination"`
	Data       []domain.Submissions `json:"data"`
}
