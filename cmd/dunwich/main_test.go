package main

import (
	"fmt"
	"syscall"
	"testing"
	"time"
)

func TestHandleSHutDownShouldCatchKill(t *testing.T) {

	done := make(chan bool, 1)

	go func() {
		handleShutdown()
		done <- true
	}()

	syscall.Kill(syscall.Getpid(), syscall.SIGKILL)

	select {
	case <-done:
		return
	case <-time.After(3 * time.Second):
		fmt.Println("timeout 2")
		t.Errorf("failed to handle %s", syscall.SIGKILL)
	}
}
