package cmd

import (
	appConfiguration "boiler-plate/app/appconf"
	"boiler-plate/internal/base/handler"
	urlHandler "boiler-plate/internal/url/handler"
	URLRepo "boiler-plate/internal/url/repository"
	"boiler-plate/internal/url/repository/redisser"
	URLService "boiler-plate/internal/url/service"
	"boiler-plate/pkg/db"
	"boiler-plate/pkg/httpclient"
	"boiler-plate/pkg/migration"
	"boiler-plate/pkg/xvalidator"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	appConf        *appConfiguration.Config
	baseHandler    *handler.BaseHTTPHandler
	URLHandler     *urlHandler.HTTPHandler
	sqlClientRepo  *db.SQLClientRepository
	validate       *validator.Validate
	httpClient     httpclient.Client
	xvalidate      *xvalidator.Validator
	urlGrpcHandler *urlHandler.GRPCHandler
	redisClient    redisser.RedisClient
)

func initHttpclient() {
	httpClientFactory := httpclient.New()
	httpClient = httpClientFactory.CreateClient()
}

//func initGRPC() {
//	//grpcHandler = handler.NewGRPCHandler(pb.UnimplementedGreeterServer{})
//	grpcHandler = &handler.GRPCHandler{
//		UnimplementedGreeterServer: pb.UnimplementedGreeterServer{},
//	}
//	usersGrpcHandler = &tempHandler.GRPCHandler{
//		UnimplementedServiceServer: users.UnimplementedServiceServer{},
//		UsersService:               nil,
//	}
//}

func initHTTP() {
	initValidator()
	appConf = appConfiguration.InitAppConfig(xvalidate)
	initInfrastructure(appConf)

	// appConf.MysqlTZ = postgresClientRepo.TZ

	baseHandler = handler.NewBaseHTTPHandler(sqlClientRepo.DB, appConf, sqlClientRepo, httpClient)

	URLRepo := URLRepo.NewRepository(sqlClientRepo.DB, sqlClientRepo)

	URLService := URLService.NewService(appConf, URLRepo, sqlClientRepo.DB, redisClient, validate)
	urlGrpcHandler = urlHandler.NewGRPCHandler(URLService)
	URLHandler = urlHandler.NewHTTPHandler(baseHandler, urlGrpcHandler, URLService)

}

func initInfrastructure(config *appConfiguration.Config) {
	initSQL(config)
	initRedis(config)
	initHttpclient()
	initLog()
}
func initValidator() {
	validate = validator.New()
	xvalidate = xvalidator.NewValidator()
}

func initRedis(config *appConfiguration.Config) {
	var ctx = context.TODO()
	rdb, _ := strconv.Atoi(config.RedisConfig.Redisdatabase)
	r := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.RedisConfig.Redishost, config.RedisConfig.Redisport),
		Password: config.RedisConfig.Redispassword,
		DB:       rdb,
	})

	err := r.Ping(ctx).Err()
	if err != nil {
		log.Fatal(err)
	}

	redisClient = redisser.NewRedisClient(r)
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
