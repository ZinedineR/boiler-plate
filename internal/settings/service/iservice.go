package service

import (
	"boiler-plate/internal/settings/domain"
	"context"
)

type Service interface {
	FindSettings(ctx context.Context) (*domain.MainTable, error)
}
