package appconf

import (
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	AppEnv       string
	AppDebug     string
	AppVersion   string   `validate:"required,startswith=v,alphanum" name:"APP_VERSION"`
	AppName      string   `validate:"required" name:"APP_NAME"`
	HttpPort     string   `validate:"required,number" name:"HTTP_PORT"`
	AllowOrigins []string `validate:"required" name:"ALLOW_ORIGINS"`
	AllowMethods []string `validate:"required" name:"ALLOW_METHODS"`
	AllowHeaders []string `validate:"required" name:"ALLOW_HEADERS"`
	FilePath     string   `validate:"required" name:"FILE_PATH"`
	FileMaxSize  int      `validate:"required,number" name:"FILE_MAX_SIZE"`
}

func AppConfigInit() *AppConfig {
	maxSize := os.Getenv("FILE_MAX_SIZE")
	maxSizeInt, err := strconv.Atoi(maxSize)
	if err != nil {
		logrus.Panicf("FILE_MAX_SIZE must be int")
	}
	return &AppConfig{
		AppEnv:       os.Getenv("APP_ENV"),
		AppDebug:     os.Getenv("APP_DEBUG"),
		AppVersion:   os.Getenv("APP_VERSION"),
		AppName:      os.Getenv("APP_NAME"),
		HttpPort:     os.Getenv("HTTP_PORT"),
		AllowOrigins: strings.Split(os.Getenv("ALLOW_ORIGINS"), ","),
		AllowMethods: strings.Split(os.Getenv("ALLOW_METHODS"), ","),
		AllowHeaders: strings.Split(os.Getenv("ALLOW_HEADERS"), ","),
		FilePath:     os.Getenv("FILE_PATH"),
		FileMaxSize:  maxSizeInt,
	}
}
