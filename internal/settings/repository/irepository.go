package repository

import (
	"boiler-plate/internal/settings/domain"
	"context"
)

type Repository interface {
	FindSettings(ctx context.Context) (*domain.MainTable, error)
}
