package config

import (
	"fmt"
	"net/netip"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const envPrefix = "DUNWICH"

type Config struct {
	EnableProfilling bool
	Logger           struct {
		Level log.Level
	}
	Memberlist struct {
		JoinNodes []netip.AddrPort
		Port      int
	}
}

func NewConfig() Config {
	cfg := Config{}
	cfg.Logger.Level = log.InfoLevel
	return cfg
}

func createEnvVars(name string) []string {
	return []string{fmt.Sprintf("%s_%s", envPrefix, strings.ToUpper(name))}
}

func BuildFlags(dest *Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "log",
			Aliases: []string{"l"},
			Usage:   "The `LEVEL` to log at (INFO,DEBUG,TRACE,...)",
			EnvVars: createEnvVars("log"),
			Value:   "INFO",
			Action: func(ctx *cli.Context, s string) error {
				lvl, err := log.ParseLevel(s)
				if err != nil {
					log.WithFields(log.Fields{
						"error": err,
					}).Fatal("failed to parse log level")
				}

				dest.Logger.Level = lvl
				return nil
			},
		},
		&cli.IntFlag{
			Name:        "cport",
			Aliases:     []string{"c"},
			Usage:       "The `PORT` to listen on for gossip communications",
			Destination: &dest.Memberlist.Port,
			EnvVars:     createEnvVars("cluster_port"),
			Value:       7946,
		},
		&cli.StringSliceFlag{
			Name:    "join",
			Aliases: []string{"j"},
			Usage:   "The optional existing `NODES` in address:port format to join to separated by commas",
			EnvVars: createEnvVars("join"),
			Action: func(ctx *cli.Context, s []string) error {
				for _, v := range s {
					v := v
					parsed, err := netip.ParseAddrPort(v)
					if err != nil {
						return err
					}
					dest.Memberlist.JoinNodes = append(dest.Memberlist.JoinNodes, parsed)
				}
				return nil
			},
		},
		&cli.BoolFlag{
			Name:        "pprof",
			Usage:       "Enable the HTTP pprof server on :6060",
			Destination: &dest.EnableProfilling,
			EnvVars:     createEnvVars("pprof"),
			Value:       false,
		},
	}
}
