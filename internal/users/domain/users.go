package domain

import (
	"boiler-plate/internal/wallet/domain"
	"os"
	"time"
)

const (
	UsersTableName = "users"
)

type Users struct {
	ID        int            `gorm:"primaryKey;not null;autoIncrement" json:"id"`
	Name      string         `json:"name"`
	Email     string         `validate:"required,gt=2" json:"email,omitempty"`
	Password  string         `json:"password,omitempty"`
	WalletId  int            `gorm:"not null" json:"wallet_id"`
	Wallet    *domain.Wallet `gorm:"foreignKey:WalletId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"wallet,omitempty"`
	CreatedAt *time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}

func (model *Users) TableName() string {
	return os.Getenv("DB_PREFIX") + UsersTableName
}
