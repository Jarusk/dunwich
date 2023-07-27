package config

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestNewConfigLevelIsInfo(t *testing.T) {
	cfg := NewConfig()

	if cfg.Logger.Level != log.InfoLevel {
		t.Errorf("expected %s, got %s", log.InfoLevel.String(), cfg.Logger.Level.String())
	}
}
