package domain

import (
	"os"
	"time"
)

const (
	URLTableName = "url"
)

type URL struct {
	ID        int        `gorm:"primaryKey;not null;autoIncrement" json:"id" example:"1"`
	Longurl   string     `validate:"required" json:"longurl,omitempty" example:"http://google.com"`
	Shorturl  string     `json:"shortutl,omitempty" example:"http://google.com"`
	CreatedAt *time.Time `gorm:"autoCreateTime" json:"created_at" example:"2024-07-19T10:57:42.454071+07:00"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime" json:"updated_at" example:"2024-07-19T10:57:42.454071+07:00"`
}

func (model *URL) TableName() string {
	return os.Getenv("DB_PREFIX") + URLTableName
}
