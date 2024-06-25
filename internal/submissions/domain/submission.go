package domain

import (
	"boiler-plate/internal/users/domain"
	"encoding/json"
	"os"
	"time"
)

const (
	SubmissionsTableName = "submissions"
)

type Submissions struct {
	ID             int             `gorm:"primaryKey;not null;autoIncrement" json:"id"`
	UserId         int             `json:"user_id"`
	User           *domain.Users   `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
	Answers        json.RawMessage `gorm:"type:json" json:"answers,omitempty"`
	RiskScore      int             `json:"risk_score"`
	RiskCategory   string          `json:"risk_category"`
	RiskDefinition string          `json:"risk_definition"`
	CreatedAt      *time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      *time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

type Question struct {
	ID       int      `json:"id"`
	Question string   `json:"question"`
	Options  []Option `json:"options"`
}

type Option struct {
	Answer string `json:"answer"`
	Weight int    `json:"weight"`
}

type RiskProfile struct {
	MinScore   int    `json:"min_score"`
	MaxScore   int    `json:"max_score"`
	Category   string `json:"category"`
	Definition string `json:"definition"`
}

func (model *Submissions) TableName() string {
	return os.Getenv("DB_PREFIX") + SubmissionsTableName
}
