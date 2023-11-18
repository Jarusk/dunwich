package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Jarusk/dunwich/pkg/carriers/tcp"
	"github.com/Jarusk/dunwich/pkg/cluster"
	"github.com/Jarusk/dunwich/pkg/config"
	"github.com/carlmjohnson/versioninfo"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"net/http"
	_ "net/http/pprof"
)

// init only needs to setup logrus currently
func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)

	var programLevel = new(slog.LevelVar) // Info by default

	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(h))

	programLevel.Set(slog.LevelInfo)

	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()
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
		Commands: []*cli.Command{
			{
				Name:    "direct",
				Aliases: []string{"d"},
				Usage:   "direct (client -> server) mode",
				Subcommands: []*cli.Command{
					{
						Name:    "client",
						Aliases: []string{"c"},
						Usage:   "start in client mode",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "server",
								Aliases: []string{"s"},
								Value:   "localhost:4365",
								Usage:   "the server to connect to, `IP:PORT`",
							},
						},
						Action: func(cCtx *cli.Context) error {
							serverAddr := cCtx.String("server")

							slog.Info("creating new client connection", "addr", serverAddr)

							var wg sync.WaitGroup
							slog.Debug("creating new clients")

							for i := 0; i < 5; i++ {
								wg.Add(1)

								go func() {
									defer wg.Done()
									c := tcp.NewTcpClient()
									slog.Debug("starting new client", "server", serverAddr)
									c.StartClient(serverAddr)
								}()
							}

							slog.Debug("waiting for signal")
							handleShutdown()

							slog.Debug("caught signal, shutting down")
							//c.Shutdown()

							slog.Debug("exiting")

							return nil
						},
					},
					{
						Name:    "server",
						Aliases: []string{"s"},
						Usage:   "start in server mode",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "listen",
								Aliases: []string{"l"},
								Value:   "0.0.0.0:4365",
								Usage:   "the listen address for the server, `IP:PORT`",
							},
						},
						Action: func(cCtx *cli.Context) error {
							slog.Info("new server listening", slog.String("addr", cCtx.String("listen")))

							slog.Debug("creating new TcpServer")
							s := tcp.NewTcpServer()
							slog.Debug("starting to accept connection")
							s.StartServer(cCtx.String("listen"))

							slog.Debug("waiting for signal")
							handleShutdown()
							slog.Debug("caught signal to shutdown, closing server")

							s.Shutdown()

							slog.Debug("finished, exiting")

							return nil
						},
					},
				},
			},
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
