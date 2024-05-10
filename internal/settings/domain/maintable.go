package domain

import (
	"os"
	"time"
)

const (
	SettingsTableName = "settings"
)

type MainTable struct {
	ID                        int        `gorm:"primaryKey;not null;autoIncrement" json:"id"`
	Currency                  string     `json:"currency,omitempty"`
	TaxFee                    float64    `json:"tax_fee,omitempty"`
	ReminderTaxProfileExpired int        `json:"reminder_tax_profile_expired,omitempty"`
	ValidAccountExpired       int        `json:"valid_account_expired,omitempty"`
	AccountExpiredPeriod      string     `json:"account_expired_period,omitempty" gorm:"default:days" validate:"eq=days|eq=months|eq=years"`
	LogoImage                 string     `json:"logo_image,omitempty"`
	Favicon                   string     `json:"favicon,omitempty"`
	PasswordLength            int        `json:"password_length,omitempty"`
	PasswordInvalid           int        `json:"password_invalid,omitempty" gorm:"default:3"`
	PasswordExpirationCount   int        `json:"password_expiration_count,omitempty"`
	PasswordExpiredPeriod     string     `json:"password_expired_period,omitempty" gorm:"default:days" validate:"eq=days|eq=months|eq=years"`
	ExpirationReminderDay     int        `json:"expiration_reminder_day,omitempty"`
	PasswordCycle             int        `json:"password_cycle,omitempty"`
	ComplexityNumeric         bool       `json:"complexity_numeric" gorm:"default:false"`
	ComplexityAlphabet        bool       `json:"complexity_alphabet" gorm:"default:false"`
	ComplexityUppercase       bool       `json:"complexity_uppercase" gorm:"default:false"`
	ComplexitySymbol          bool       `json:"complexity_symbol" gorm:"default:false"`
	UpdatedAt                 *time.Time `json:"updated_at"`
}

func (model *MainTable) TableName() string {
	return os.Getenv("DB_PREFIX") + SettingsTableName
}
