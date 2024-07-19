package repository

import (
	"boiler-plate/internal/url/domain"
	"context"
	"gorm.io/gorm"
)

type UsersRepository interface {
	Create(ctx context.Context, tx *gorm.DB, model *domain.URL) error
	Detail(ctx context.Context, tx *gorm.DB, longurl string) (*domain.URL, error)
}
