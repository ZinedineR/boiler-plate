package appconf

import (
	"os"
	"strconv"
)

type RedisConfig struct {
	Redishost     string `validate:"required" name:"REDIS_HOST"`
	Redisport     int    `validate:"required" name:"REDIS_PORT"`
	Redisdatabase string `validate:"required" name:"REDIS_DATABASE"`
	Redispassword string `name:"REDIS_PASSWORD"`
}

func RedisConfigInit() *RedisConfig {
	port, _ := strconv.Atoi(os.Getenv("REDIS_PORT"))
	return &RedisConfig{
		Redishost:     os.Getenv("REDIS_HOST"),
		Redisport:     port,
		Redisdatabase: os.Getenv("REDIS_DATABASE"),
		Redispassword: os.Getenv("REDIS_PASSWORD"),
	}
}
