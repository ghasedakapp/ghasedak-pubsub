package pkg

import (
	"os"
	"os/signal"
	"syscall"
)

func Wait() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		GetLogger().Debug(sig)
		done <- true
	}()

	GetLogger().Debug("awaiting signal")
	<-done
	GetLogger().Debug("exiting")
}
