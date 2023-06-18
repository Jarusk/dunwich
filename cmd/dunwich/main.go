package main

import (
	"fmt"
	"os"

	"github.com/Jarusk/dunwich/pkg/config"
	log "github.com/sirupsen/logrus"
)

// init only needs to setup logrus currently
func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

func setLogLevel(cfg *config.Config) {
	lvl, err := log.ParseLevel(cfg.Logger.Level)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("failed to parse log level")
	}

	log.SetLevel(lvl)
	log.Infof("set log level to '%s'", lvl.String())
}

func main() {
	log.Info("starting Dunwich")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("failed to load config")
	}

	log.WithFields(log.Fields{
		"config": fmt.Sprintf("%+v", cfg),
	}).Info("loaded config")

	setLogLevel(cfg)
}
