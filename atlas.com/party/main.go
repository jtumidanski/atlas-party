package main

import (
	"atlas-party/logger"
	"atlas-party/party"
	"atlas-party/rest"
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	l := logger.CreateLogger()
	l.Infoln("Starting main service.")

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	//consumers.CreateEventConsumers(l, ctx, wg)

	rest.CreateService(l, ctx, wg, "/ms/party", party.InitResource)

	party.GetRegistry().Create(1, 1, 1)
	party.GetRegistry().Create(1, 1, 2)

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Infof("Initiating shutdown with signal %s.", sig)
	cancel()
	wg.Wait()
	l.Infoln("Service shutdown.")
}
