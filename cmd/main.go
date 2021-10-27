package main

import (
	"os"

	cmd "github.com/ful09003/proxl/cmd/realcmd"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := cmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
