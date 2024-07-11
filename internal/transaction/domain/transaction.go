package domain

import (
	"boiler-plate/internal/users/domain"
	"os"
	"time"
)

const (
	TransactionTableName = "transaction"
)

type Transaction struct {
	ID              int           `gorm:"primaryKey;not null;autoIncrement" json:"id"`
	Type            string        `json:"type"`
	Amount          float64       `json:"amount"`
	Message         string        `json:"message"`
	UserId          int           `json:"user_id"`
	Users           *domain.Users `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"users,omitempty"`
	TransactionTime *time.Time    `gorm:"autoCreateTime" json:"transaction_time"`
}

func (model *Transaction) TableName() string {
	return os.Getenv("DB_PREFIX") + TransactionTableName
}
