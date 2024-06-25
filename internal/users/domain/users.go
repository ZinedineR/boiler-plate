package domain

import (
	"os"
	"time"
)

const (
	UsersTableName = "users"
)

type Users struct {
	ID        int        `gorm:"primaryKey;not null;autoIncrement" json:"id"`
	Email     string     `validate:"required,gt=2" json:"email,omitempty"`
	Password  string     `json:"password,omitempty"`
	CreatedAt *time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type UserResponse struct {
	ID             int        `gorm:"primaryKey;not null;autoIncrement" json:"id"`
	Email          string     `validate:"required,gt=2" json:"email,omitempty"`
	Password       string     `json:"password,omitempty"`
	RiskScore      int        `json:"risk_score"`
	RiskCategory   string     `json:"risk_category"`
	RiskDefinition string     `json:"risk_definition"`
	CreatedAt      *time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      *time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (model *UserResponse) DeclareRiskProfile(riskscore int, riskcategory, riskdefinition string) {
	model.RiskScore = riskscore
	model.RiskCategory = riskcategory
	model.RiskDefinition = riskdefinition
}

func (model *Users) TableName() string {
	return os.Getenv("DB_PREFIX") + UsersTableName
}
