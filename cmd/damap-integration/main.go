package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tuwien-csd/damap-integration"
)

func main() {
	if len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		print(damap_integration.ConfigHelp())
		return
	}

	integration, err := damap_integration.New(
		damap_integration.WithEnv(os.Environ()),
		damap_integration.WithCLI(os.Args),
		damap_integration.WithLogger(slog.New(slog.NewTextHandler(os.Stderr, nil))),
	)
	if err != nil {
		panic(err)
	}

	runningIntegration, err := integration.Start()
	if err != nil {
		panic(err)
	}

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	shutdown := func() {
		shutdownContext, c := context.WithTimeout(context.Background(), 30*time.Second)
		defer c()
		if err := runningIntegration.Stop(shutdownContext); err != nil {
			panic(err)
		}
	}

	select {
	case <-sigChannel:
		shutdown()
	case <-runningIntegration.OnStop():
		shutdown()
	}
}
