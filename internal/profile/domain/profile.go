package domain

import (
	"os"
	"time"
)

const (
	ProfileTableName = "profile"
)

type Profile struct {
	ID        int        `gorm:"primaryKey;not null;autoIncrement" json:"id"`
	Profile   string     `json:"profile,omitempty"`
	Password  string     `json:"password,omitempty"`
	CreatedAt *time.Time `json:"created_at"`
}

func (model *Profile) TableName() string {
	return os.Getenv("DB_PREFIX") + ProfileTableName
}
