package migration

import (
	submissionDomain "boiler-plate/internal/submissions/domain"
	"boiler-plate/internal/users/domain"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sort"
)

func Initmigrate(db *gorm.DB) {

	sqlDB, err := db.DB()
	if err != nil {
		defer sqlDB.Close()
	}

	executePendingMigrations(db)

	// Migrate rest of the models
	logrus.Println(fmt.Println("AutoMigrate Model [table_name]"))
	db.AutoMigrate(&domain.Users{})
	logrus.Infoln(fmt.Println("  TableModel [" +
		(&domain.Users{}).TableName() + "]"))
	db.AutoMigrate(&submissionDomain.Submissions{})
	logrus.Infoln(fmt.Println("  TableModel [" +
		(&submissionDomain.Submissions{}).TableName() + "]"))
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

// func registerMigration(id string, fm mFunc) {
// 	migrations[id] = fm
// }
