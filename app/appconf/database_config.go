package appconf

import (
	"os"
	"strconv"
)

type DatabaseConfig struct {
	Dbservice  string `validate:"required,eq=postgres|eq=mysql|eq=sqlserver" name:"DB_CONNECTION"`
	Dbhost     string `validate:"required" name:"DB_HOST"`
	Dbport     int    `validate:"required" name:"DB_PORT"`
	Dbdatabase string `validate:"required" name:"DB_DATABASE"`
	Dbuser     string `validate:"required" name:"DB_USERNAME"`
	Dbpassword string `validate:"required" name:"DB_PASSWORD"`
	Dbprefix   string `validate:"required" name:"DB_PREFIX"`
}

func DatabaseConfigInit() *DatabaseConfig {
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	return &DatabaseConfig{
		Dbservice:  os.Getenv("DB_CONNECTION"),
		Dbhost:     os.Getenv("DB_HOST"),
		Dbport:     port,
		Dbdatabase: os.Getenv("DB_DATABASE"),
		Dbuser:     os.Getenv("DB_USERNAME"),
		Dbpassword: os.Getenv("DB_PASSWORD"),
		Dbprefix:   os.Getenv("DB_PREFIX"),
	}
}
