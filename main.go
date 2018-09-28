package main

import (
	log "github.com/sirupsen/logrus"
	"github.frg.tech/cloud/fanplane/cmd"
	"os"
)

func main() {
	if err := cmd.FanplaneCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(-1)
	}
}
