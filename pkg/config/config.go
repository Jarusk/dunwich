package config

import (
	"github.com/kkyr/fig"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Logger struct {
		Level string `fig:"level" default:"info"`
	}
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := fig.Load(
		&cfg,
		fig.IgnoreFile(),
		fig.UseEnv("DUNWICH"),
	)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to load config")
		return nil, err
	}
	return &cfg, nil
}
