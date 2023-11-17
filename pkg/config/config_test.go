package config

import (
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestCreateEnvVars(t *testing.T) {
	cases := map[string]struct {
		input       string
		expected    []string
		shouldMatch bool
	}{
		"shouldSucceedWithSingle": {
			input:       "log",
			expected:    []string{"DUNWICH_LOG"},
			shouldMatch: true,
		},
		"shouldSucceedWithUnderscore": {
			input:       "stuff_things",
			expected:    []string{"DUNWICH_STUFF_THINGS"},
			shouldMatch: true,
		},
		"shoudlFailNonmatching": {
			input:       "fail",
			expected:    []string{"DUNWICH_SUCCEED"},
			shouldMatch: false,
		},
	}

	t.Parallel()
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			res := createEnvVars(tc.input)
			match := reflect.DeepEqual(res, tc.expected)

			if match && !tc.shouldMatch || !match && tc.shouldMatch {
				t.Errorf("expected %v, have %v", tc.expected, res)
			}
		})
	}
}

func TestNewConfigLevelIsInfo(t *testing.T) {
	t.Parallel()

	cfg := NewConfig()

	if cfg.Logger.Level != log.InfoLevel {
		t.Errorf("expected %s, got %s", log.InfoLevel.String(), cfg.Logger.Level.String())
	}
}
