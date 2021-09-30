package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func checkForSignals(errCollector chan error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	errCollector <- fmt.Errorf("%s", <-c)
}
