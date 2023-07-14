package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Jarusk/dunwich/pkg/cluster"
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

func handleShutdown() {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	log.Info("finished setup")

	caught := <-stop

	log.WithFields(log.Fields{
		"signal": caught.String(),
	}).Info("caught signal, exiting.")
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

	// externalize
	joinNode := flag.String("join", "", "the address:port to point to")
	flag.Parse()

	bootstrapNodes := make([]string, 0, 1)

	if *joinNode != "" {
		bootstrapNodes = append(bootstrapNodes, *joinNode)
	}

	cluster := cluster.Cluster{}
	cluster.JoinCluster(cfg.Memberlist.Port, bootstrapNodes)

	handleShutdown()
}
