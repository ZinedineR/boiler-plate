package repository

import (
	"boiler-plate/internal/profile/domain"
	"context"
	"gorm.io/gorm"
)

type ProfileRepository interface {
	Create(ctx context.Context, tx *gorm.DB, model *domain.Profile) error
	Update(ctx context.Context, tx *gorm.DB, id int, model *domain.Profile) error
	Find(ctx context.Context, tx *gorm.DB) (*[]domain.Profile, error)
	Detail(ctx context.Context, tx *gorm.DB, id int) (*domain.Profile, error)
	Auth(ctx context.Context, tx *gorm.DB, profile, password string) (*domain.Profile, error)
	Delete(ctx context.Context, tx *gorm.DB, key int) error
}
