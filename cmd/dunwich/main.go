package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/Jarusk/dunwich/pkg/cluster"
	"github.com/Jarusk/dunwich/pkg/config"
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

func handleShutdown(stop chan os.Signal) {
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	log.Info("finished setup")

	caught := <-stop

	log.WithFields(log.Fields{
		"signal": caught.String(),
	}).Info("caught signal, exiting.")
}

func main() {
	log.Trace("starting Dunwich")

	cfg := config.NewConfig()

	info, _ := debug.ReadBuildInfo()

	app := &cli.App{
		Name:                 "dunwich",
		Usage:                "A next generation network load tester",
		Version:              info.Main.Version,
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

			handleShutdown(make(chan os.Signal, 1))
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
