package repository

import (
	"boiler-plate/internal/users/domain"
	"boiler-plate/pkg/db"
	"context"
	"gorm.io/gorm"
)

type UsersRepository interface {
	Create(ctx context.Context, tx *gorm.DB, model *domain.Users) error
	Update(ctx context.Context, tx *gorm.DB, id int, model *domain.Users) error
	Find(ctx context.Context, tx *gorm.DB, limit, page int) (*[]domain.Users, *db.Paginate, error)
	Detail(ctx context.Context, tx *gorm.DB, id int) (*domain.Users, error)
	Delete(ctx context.Context, tx *gorm.DB, key int) error
}
