package appconf

import (
	"boiler-plate/pkg/xvalidator"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

type Config struct {
	AppEnvConfig   *AppConfig
	DatabaseConfig *DatabaseConfig
	RedisConfig    *RedisConfig
}

func (c Config) IsStaging() bool {
	return c.AppEnvConfig.AppEnv != "production"
}

func (c Config) IsProd() bool {
	return c.AppEnvConfig.AppEnv == "production"
}

func (c Config) IsDebug() bool {
	return c.AppEnvConfig.AppDebug == "True"
}

func InitAppConfig(validate *xvalidator.Validator) *Config {

	c := Config{
		AppEnvConfig:   AppConfigInit(),
		DatabaseConfig: DatabaseConfigInit(),
		RedisConfig:    RedisConfigInit(),
	}

	errs := validate.Struct(c)
	if errs != nil {
		for _, e := range errs {
			logrus.Error(fmt.Sprintf("Failed to load env: %s", e))
		}
		os.Exit(1)
	}
	return &c

}
