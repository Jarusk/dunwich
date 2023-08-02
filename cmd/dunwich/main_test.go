package main

import (
	"os"
	"syscall"
	"testing"
	"time"
)

func TestHandleShutDownShouldCatch(t *testing.T) {
	tests := map[string]struct {
		sig         syscall.Signal
		shouldCatch bool
	}{
		"shouldCatchSIGINT": {
			sig:         syscall.SIGINT,
			shouldCatch: true,
		},
		"shouldCatchSIGTERM": {
			sig:         syscall.SIGTERM,
			shouldCatch: true,
		},
		"shouldCatchSIGHUP": {
			sig:         syscall.SIGHUP,
			shouldCatch: true,
		},
		"shouldNotCatchSIGPIPE": {
			sig:         syscall.SIGPIPE,
			shouldCatch: false,
		},
	}

	t.Parallel()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			done := make(chan bool, 1)
			testPid := make(chan int, 1)

			go func() {
				testPid <- os.Getpid()
				handleShutdown()
				done <- true
			}()

			p, err := os.FindProcess(<-testPid)
			if err != nil {
				panic(err)
			}

			// Give enough time for the signal handlers to be registered
			// before sending the signal
			time.Sleep(100 * time.Millisecond)
			if err = p.Signal(tc.sig); err != nil {
				t.Errorf("failed to send signal to child: %s", err)
			}

			select {
			case <-done:
				if !tc.shouldCatch {
					t.Errorf("handled but should not have: %+v", tc)
				}
				return
			case <-time.After(1 * time.Second):
				if tc.shouldCatch {
					t.Errorf("failed to handle after 1 second: %+v", tc)
				}
			}
		})
	}
}
