package main

import (
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestHandleShutDownShouldCatch(t *testing.T) {
	tests := []syscall.Signal{
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGHUP,
	}
	t.Parallel()
	for _, tc := range tests {
		t.Run(tc.String(), func(t *testing.T) {
			done := make(chan bool, 1)
			notification := make(chan os.Signal, 1)

			go func() {
				handleShutdown(notification)
				done <- true
			}()

			notification <- tc

			select {
			case <-done:
				return
			case <-time.After(3 * time.Second):
				fmt.Println("timeout 2")
				t.Errorf("failed to handle %s", tc)
			}
		})
	}
}
