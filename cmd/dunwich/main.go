package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Jarusk/dunwich/pkg/cluster"
	"github.com/Jarusk/dunwich/pkg/config"
	"github.com/carlmjohnson/versioninfo"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// init only needs to setup logrus currently
func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

func handleShutdown() {
	log.Info("finished setup")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	caught := <-stop

	log.WithFields(log.Fields{
		"signal": caught.String(),
	}).Info("caught signal, handling shutdown.")

	// Cleanup logic goes here.
}

func run(args []string) {

	log.Tracef("recevied args %v", args)

	cfg := config.NewConfig()

	app := &cli.App{
		Name:                 "dunwich",
		Usage:                "A next generation network load tester",
		Version:              versioninfo.Short(),
		EnableBashCompletion: true,
		Flags:                config.BuildFlags(&cfg),
		Suggest:              true,
		Action: func(ctx *cli.Context) error {
			log.SetLevel(cfg.Logger.Level)
			log.Infof("set log level to '%s'", cfg.Logger.Level.String())

			log.WithFields(log.Fields{
				"config": fmt.Sprintf("%+v", cfg),
			}).Debug("loaded config")

			cluster := cluster.Cluster{}
			cluster.JoinCluster(cfg.Memberlist.Port, cfg.Memberlist.JoinNodes)

			handleShutdown()

			return nil
		},
	}

	if err := app.Run(args); err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.Trace("starting Dunwich")
	run(os.Args)
}
