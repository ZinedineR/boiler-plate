package cmd

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
	appConfiguration "ms-batch/app/appconf"
	"ms-batch/internal/base/handler"
	TemplateService "ms-batch/internal/template/service"
	"ms-batch/pkg/db"
	"ms-batch/pkg/httpclient"
	"os"

	"github.com/sirupsen/logrus"
	tempHandler "ms-batch/internal/template/handler"
	templateRepo "ms-batch/internal/template/repository"
)

var (
	appConf         *appConfiguration.Config
	baseHandler     *handler.BaseHTTPHandler
	templateHandler *tempHandler.HTTPHandler
	mongoClientRepo *db.MongoDBClientRepository
	validate        *validator.Validate
	httpClient      httpclient.Client
)

func initMongoDB() {
	mongoClientRepo, _ = db.NewMongoDBRepository("", "", "", 0)
}

func initHttpclient() {
	httpClientFactory := httpclient.New()
	httpClient = httpClientFactory.CreateClient()
}

func initHTTP() {
	appConf = appConfiguration.InitAppConfig()
	initInfrastructure()

	//appConf.MysqlTZ = postgresClientRepo.TZ

	baseHandler = handler.NewBaseHTTPHandler(mongoClientRepo.Client, appConf, mongoClientRepo, httpClient)

	templateRepo := templateRepo.NewRepository(mongoClientRepo.DB, mongoClientRepo)
	templateService := TemplateService.NewService(templateRepo, validate)
	templateHandler = tempHandler.NewHTTPHandler(baseHandler, templateService)
}

func initInfrastructure() {
	//initPostgreSQL()
	initHttpclient()
	initMongoDB()
	initLog()
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
