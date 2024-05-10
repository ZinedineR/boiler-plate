package migration

import (
	"boiler-plate/internal/settings/domain"
	"encoding/base64"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func Initmigrate(db *gorm.DB) {

	sqlDB, err := db.DB()
	if err != nil {
		defer sqlDB.Close()
	}

	executePendingMigrations(db)

	// Migrate rest of the models
	logrus.Println(fmt.Println("AutoMigrate Model [table_name]"))
	db.AutoMigrate(&domain.MainTable{})
	logrus.Infoln(fmt.Println("  TableModel [" +
		(&domain.MainTable{}).TableName() + "]"))
	// Check if the Category table is empty.
	// Check if the Category table is empty.
	var settingsCount int64
	db.Model(&domain.MainTable{}).Count(&settingsCount)

	if settingsCount == 0 {
		// The Category table is empty, so seed data.
		seedDatabase(db)
	}
}

func executePendingMigrations(db *gorm.DB) {
	db.AutoMigrate(&MigrationHistoryModel{})
	lastMigration := MigrationHistoryModel{}
	skipMigration := db.Order("migration_id desc").Limit(1).Find(&lastMigration).RowsAffected > 0

	// skip to last migration
	keys := make([]string, 0, len(migrations))
	for k := range migrations {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// run all migrations in one transaction
	if len(migrations) == 0 {
		logrus.Infoln(fmt.Print("No pending migrations"))
	} else {
		db.Transaction(func(tx *gorm.DB) error {
			for _, k := range keys {
				if skipMigration {
					if k == lastMigration.MigrationID {
						skipMigration = false
					}
				} else {
					logrus.Infoln(fmt.Sprintf("  " + k))
					tx.Transaction(func(subTx *gorm.DB) error {
						// run migration update
						checkError(migrations[k](subTx))
						// insert migration id into history
						checkError(subTx.Create(MigrationHistoryModel{MigrationID: k}).Error)
						return nil
					})
				}
			}
			return nil
		})
	}
}

type mFunc func(tx *gorm.DB) error

var migrations = make(map[string]mFunc)

// MigrationHistoryModel model migration
type MigrationHistoryModel struct {
	MigrationID string `gorm:"primaryKey"`
}

// TableName name of migration table
func (model *MigrationHistoryModel) TableName() string {
	return "migration_history"
}

func checkError(err error) {
	if err != nil {
		logrus.Infoln(fmt.Print(err.Error()))
		panic(err)
	}
}

func seedDatabase(db *gorm.DB) {
	// Create and insert 1-2 rows of data for the Article table.
	timeFormat := "02-January-2006-15-4-5"

	newFilename := "default-" + time.Now().Format(timeFormat) + ".png"
	Settings := &domain.MainTable{
		ID:                        0,
		Currency:                  "IDR",
		TaxFee:                    1.0,
		ReminderTaxProfileExpired: 1,
		ValidAccountExpired:       30,
		LogoImage:                 newFilename,
		Favicon:                   newFilename,
		PasswordLength:            8,
		PasswordExpirationCount:   12,
		ExpirationReminderDay:     1,
		ComplexityNumeric:         true,
	}
	data, err := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAYAAACNMs+9AAAAFUlEQVR42mNk+M9Qz0AEYBxVSF+FAAhKDveksOjmAAAAAElFTkSuQmCC")
	if err != nil {
		return
	}
	filePath := os.Getenv("FILE_PATH")
	dst, err := os.Create(filePath + newFilename)
	defer dst.Close()
	if _, err = dst.Write(data); err != nil {
		return
	}

	db.Create(&Settings)

}

// func registerMigration(id string, fm mFunc) {
// 	migrations[id] = fm
// }
