package repository

import (
	"boiler-plate/internal/submissions/domain"
	"boiler-plate/pkg/db"
	"context"
	"gorm.io/gorm"
)

type SubmissionsRepository interface {
	Create(ctx context.Context, tx *gorm.DB, model *domain.Submissions) error
	Update(ctx context.Context, tx *gorm.DB, id int, model *domain.Submissions) error
	Find(ctx context.Context, tx *gorm.DB, limit, page int) (*[]domain.Submissions, *db.Paginate, error)
	FindByUser(ctx context.Context, tx *gorm.DB, limit, page, userid int) (*[]domain.Submissions, *db.Paginate, error)
	Detail(ctx context.Context, tx *gorm.DB, id int) (*domain.Submissions, error)
	DetailByUser(ctx context.Context, tx *gorm.DB, id int) (*domain.Submissions, error)
	Delete(ctx context.Context, tx *gorm.DB, key int) error
}
