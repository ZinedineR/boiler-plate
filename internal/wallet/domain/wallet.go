package domain

import (
	"os"
	"time"
)

const (
	WalletTableName = "wallet"
)

type Wallet struct {
	ID              int        `gorm:"primaryKey;not null;autoIncrement" json:"id"`
	Balance         float64    `gorm:"default:0" json:"balance"`
	LastTransaction *time.Time `gorm:"autoUpdateTime" json:"last_transaction"`
}

func (model *Wallet) Increase(amount float64) {
	model.Balance += amount
}

func (model *Wallet) Decrease(amount float64) {
	model.Balance -= amount
}

func (model *Wallet) TableName() string {
	return os.Getenv("DB_PREFIX") + WalletTableName
}
