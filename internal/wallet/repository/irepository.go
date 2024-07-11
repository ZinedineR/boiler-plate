package repository

import (
	"boiler-plate/internal/wallet/domain"
	"context"
	"gorm.io/gorm"
)

type WalletRepository interface {
	Create(ctx context.Context, tx *gorm.DB, model *domain.Wallet) error
	Update(ctx context.Context, tx *gorm.DB, id int, model *domain.Wallet) error
}
