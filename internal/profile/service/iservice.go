package service

import (
	"boiler-plate/internal/profile/domain"
	"boiler-plate/pkg/exception"
	"context"
)

type Service interface {
	Create(
		ctx context.Context, req *domain.Profile,
	) *exception.Exception
	// Detail Service
	Detail(ctx context.Context, id string) (*domain.Profile, *exception.Exception)
	// Delete Service
	Delete(ctx context.Context, id string) *exception.Exception
	// Update Service
	Update(
		ctx context.Context, id string, profile *domain.Profile,
	) *exception.Exception
	Find(ctx context.Context) (*[]domain.Profile, *exception.Exception)
}
