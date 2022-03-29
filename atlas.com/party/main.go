package main

import (
	"atlas-party/character"
	"atlas-party/kafka"
	"atlas-party/logger"
	"atlas-party/party"
	"atlas-party/rest"
	"atlas-party/tracing"
	"context"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const serviceName = "atlas-party"
const consumerGroupId = "Party Orchestration Service"

func main() {
	l := logger.CreateLogger(serviceName)
	l.Infoln("Starting main service.")

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	tc, err := tracing.InitTracer(l)(serviceName)
	if err != nil {
		l.WithError(err).Fatal("Unable to initialize tracer.")
	}
	defer func(tc io.Closer) {
		err = tc.Close()
		if err != nil {
			l.WithError(err).Errorf("Unable to close tracer.")
		}
	}(tc)

	kafka.CreateConsumers(l, ctx, wg,
		character.StatusConsumer(consumerGroupId),
		party.CreateConsumer(consumerGroupId),
		party.ExpelConsumer(consumerGroupId),
		party.JoinConsumer(consumerGroupId),
		party.LeaveConsumer(consumerGroupId),
		party.PromoteLeaderConsumer(consumerGroupId))

	rest.CreateService(l, ctx, wg, "/ms/party", party.InitResource, character.InitResource)

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
