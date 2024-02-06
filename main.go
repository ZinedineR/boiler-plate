package main

import (
	"github.com/sirupsen/logrus"
	"ms-batch/cmd"
	"os"
)

func main() {

	if err := cmd.Execute(); err != nil {
		logrus.Errorln("error on command execution", err.Error())
		os.Exit(1)
	}
}
