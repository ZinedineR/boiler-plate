package repository

import (
	"boiler-plate/internal/transaction/domain"
	"boiler-plate/pkg/db"
	"context"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *gorm.DB, model *domain.Transaction) error
	Update(ctx context.Context, tx *gorm.DB, id int, model *domain.Transaction) error
	Find(ctx context.Context, tx *gorm.DB, limit, page, userid int) (*[]domain.Transaction, *db.Paginate, error)
	Detail(ctx context.Context, tx *gorm.DB, id int) (*domain.Transaction, error)
}
