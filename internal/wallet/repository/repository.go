package repository

import (
	"boiler-plate/internal/wallet/domain"
	baseModel "boiler-plate/pkg/db"
	"context"
	"gorm.io/gorm"
)

type Repo struct {
	db   *gorm.DB
	base *baseModel.SQLClientRepository
}

func NewRepository(db *gorm.DB, base *baseModel.SQLClientRepository) WalletRepository {
	return &Repo{db: db, base: base}
}

func (r Repo) Create(ctx context.Context, tx *gorm.DB, model *domain.Wallet) error {
	query := tx.WithContext(ctx)
	if err := query.Model(&domain.Wallet{}).Create(&model).
		Error; err != nil {
		return err
	}
	return nil
}

func (r Repo) Update(ctx context.Context, tx *gorm.DB, id int, model *domain.Wallet) error {
	query := tx.WithContext(ctx)
	if err := query.
		Model(&domain.Wallet{ID: id}).
		Updates(model).
		Error; err != nil {
		return err
	}
	return nil
}
