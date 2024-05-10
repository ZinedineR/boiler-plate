package db

import (
	appConfiguration "boiler-plate/app/appconf"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"os"
	"strconv"
	"time"
)

type SQLClientRepository struct {
	DB *gorm.DB
	TZ string
}

func NewSQLClientRepository(
	appConfig *appConfiguration.Config, config *gorm.Config,
) (*SQLClientRepository, error) {
	tz := "Asia/Jakarta"

	if config == nil {
		config = &gorm.Config{}
	}
	db, _ := gorm.Open(nil, config)
	maxIdleConn, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTION"))
	lifetimeConn, _ := time.ParseDuration(os.Getenv("DB_MAX_LIFETIME_CONNECTION"))
	var dsn string
	switch appConfig.DatabaseConfig.Dbservice {
	case "postgres", "pgsql":
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
			appConfig.DatabaseConfig.Dbhost,
			appConfig.DatabaseConfig.Dbuser,
			appConfig.DatabaseConfig.Dbpassword,
			appConfig.DatabaseConfig.Dbdatabase,
			appConfig.DatabaseConfig.Dbport, tz)
		sqlDB, err := sql.Open("pgx", dsn)
		if err != nil {
			logrus.Error(fmt.Sprintf("Cannot connect to PostgresSQL. %v", err))
			return nil, errors.Wrap(err, "Cannot connect to PostgresSQL")
		}
		sqlDB.SetMaxIdleConns(maxIdleConn)
		sqlDB.SetConnMaxLifetime(lifetimeConn)
		db, err = gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), config)
		if err != nil {
			logrus.Error(fmt.Sprintf("Cannot connect to PostgresSQL. %v", err))
			return nil, errors.Wrap(err, "Cannot connect to PostgresSQL")
		}
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			appConfig.DatabaseConfig.Dbuser,
			appConfig.DatabaseConfig.Dbpassword,
			appConfig.DatabaseConfig.Dbhost,
			appConfig.DatabaseConfig.Dbport,
			appConfig.DatabaseConfig.Dbdatabase)
		sqlDB, err := sql.Open("mysql", dsn)
		if err != nil {
			logrus.Error(fmt.Sprintf("Cannot connect to MySQL. %v", err))
			return nil, errors.Wrap(err, "Cannot connect to MySQL")
		}
		sqlDB.SetMaxIdleConns(maxIdleConn)
		sqlDB.SetConnMaxLifetime(lifetimeConn)
		db, err = gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}), config)
		if err != nil {
			logrus.Panic(fmt.Sprintf("Cannot connect to MySQL. %v", err))
			return nil, errors.Wrap(err, "Cannot connect to MySQL")
		}
		//db.Use(otelgorm.NewPlugin())
	case "sqlserver":
		var err error
		dsn = fmt.Sprintf("sqlserver://%s:%s@%s?database=%s",
			appConfig.DatabaseConfig.Dbuser,
			appConfig.DatabaseConfig.Dbpassword,
			appConfig.DatabaseConfig.Dbhost,
			appConfig.DatabaseConfig.Dbdatabase)
		db, err = gorm.Open(sqlserver.Open(dsn), config)
		sqlDB, err := db.DB()
		if err != nil {
			logrus.Errorln("failed to configure connection pool", "error", err.Error())
			os.Exit(1)
		}
		sqlDB.SetMaxIdleConns(maxIdleConn)
		sqlDB.SetConnMaxLifetime(lifetimeConn)
		if err != nil {
			logrus.Panic(fmt.Sprintf("Cannot connect to SQL Server. %v", err))
			return nil, errors.Wrap(err, "Cannot connect to SQL Server")
		}
	default:
		logrus.Errorln("unknown database driver")
		os.Exit(1)
	}

	db.Use(otelgorm.NewPlugin())

	if db == nil {
		panic("missing db")
	}

	return &SQLClientRepository{DB: db, TZ: tz}, nil

}
