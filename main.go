package main

import (
	"boiler-plate/cmd"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {

	if err := cmd.Execute(); err != nil {
		logrus.Errorln("error on command execution", err.Error())
		os.Exit(1)
	}
}
