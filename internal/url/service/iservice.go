package service

import (
	"boiler-plate/internal/url/domain"
	"boiler-plate/pkg/exception"
	"context"
)

type Service interface {
	Create(
		ctx context.Context, req *domain.URL,
	) *exception.Exception
	// Detail Service
	Detail(ctx context.Context, id string) (*domain.URL, *exception.Exception)
}
