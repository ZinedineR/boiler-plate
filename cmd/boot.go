package cmd

import (
	"boiler-plate/internal/base/handler"
	"fmt"
	"io"
	"log"
	"os"

	appConfiguration "boiler-plate/app/appconf"
	subHandler "boiler-plate/internal/submissions/handler"
	SubmissionsRepo "boiler-plate/internal/submissions/repository"
	SubmissionsService "boiler-plate/internal/submissions/service"
	tempHandler "boiler-plate/internal/users/handler"
	UsersRepo "boiler-plate/internal/users/repository"
	UsersService "boiler-plate/internal/users/service"
	"boiler-plate/pkg/db"
	"boiler-plate/pkg/httpclient"
	"boiler-plate/pkg/migration"
	"boiler-plate/pkg/xvalidator"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	appConf            *appConfiguration.Config
	baseHandler        *handler.BaseHTTPHandler
	UsersHandler       *tempHandler.HTTPHandler
	SubmissionsHandler *subHandler.HTTPHandler
	sqlClientRepo      *db.SQLClientRepository
	validate           *validator.Validate
	httpClient         httpclient.Client
	xvalidate          *xvalidator.Validator
)

func initHttpclient() {
	httpClientFactory := httpclient.New()
	httpClient = httpClientFactory.CreateClient()
}

func initHTTP() {
	initValidator()
	appConf = appConfiguration.InitAppConfig(xvalidate)
	initInfrastructure(appConf)

	// appConf.MysqlTZ = postgresClientRepo.TZ

	baseHandler = handler.NewBaseHTTPHandler(sqlClientRepo.DB, appConf, sqlClientRepo, httpClient)

	UsersRepo := UsersRepo.NewRepository(sqlClientRepo.DB, sqlClientRepo)
	SubsRepo := SubmissionsRepo.NewRepository(sqlClientRepo.DB, sqlClientRepo)

	UsersService := UsersService.NewService(appConf, UsersRepo, SubsRepo, sqlClientRepo.DB, validate)
	SubsService := SubmissionsService.NewService(appConf, SubsRepo, sqlClientRepo.DB, validate)
	UsersHandler = tempHandler.NewHTTPHandler(baseHandler, UsersService)
	SubmissionsHandler = subHandler.NewHTTPHandler(baseHandler, SubsService)

}

func initInfrastructure(config *appConfiguration.Config) {
	initSQL(config)
	initHttpclient()
	initLog()
}
func initValidator() {
	validate = validator.New()
	xvalidate = xvalidator.NewValidator()
}
func isProd() bool {
	return os.Getenv("APP_ENV") == "production"
}

func initLog() {
	lv := os.Getenv("LOG_LEVEL_DEV")
	level := logrus.InfoLevel
	switch lv {
	case "PanicLevel":
		level = logrus.PanicLevel
	case "FatalLevel":
		level = logrus.FatalLevel
	case "ErrorLevel":
		level = logrus.ErrorLevel
	case "WarnLevel":
		level = logrus.WarnLevel
	case "InfoLevel":
		level = logrus.InfoLevel
	case "DebugLevel":
		level = logrus.DebugLevel
	case "TraceLevel":
		level = logrus.TraceLevel
	default:
	}

	if isProd() {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetLevel(logrus.WarnLevel)
		logrus.SetOutput(os.Stdout)
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true})
		if lv == "" && os.Getenv("APP_DEBUG") == "True" {
			level = logrus.DebugLevel
		}
		logrus.SetLevel(level)
		// logrus.SetFormatter()
		if os.Getenv("DEV_FILE_LOG") == "True" {
			logfile, err := os.OpenFile("log/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
			if err != nil {
				fmt.Printf("error opening file : %v", err)
			}

			mw := io.MultiWriter(os.Stdout, logfile)
			logrus.SetOutput(mw)
		} else {
			logrus.SetOutput(os.Stdout)
		}
	}
}

func initSQL(config *appConfiguration.Config) {

	//var gConfig *gorm.Config
	gConfig := &gorm.Config{}
	if os.Getenv("DEV_SHOW_QUERY") == "true" {
		showQuery := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel: logger.Info,
			})
		gConfig.Logger = showQuery
	} else {
		gConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	sqlClientRepo, _ = db.NewSQLClientRepository(config, gConfig)
	if config.IsStaging() {
		migration.Initmigrate(sqlClientRepo.DB)
	}
}
