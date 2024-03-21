package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/seinshah/cvci/internal/command"
	"github.com/seinshah/cvci/internal/pkg/logger"
)

// Version will be changed at build time using ldflags.
// The new value will be determined in the CI pipeline
// using the git tag.
var Version = "development"

func main() {
	logger.SetUp()

	cmd := command.NewCommand(Version)

	ctx, cancel := context.WithCancel(context.Background())

	signalChannel := make(chan os.Signal, 1)
	exitChannel := make(chan int, 1)

	signal.Notify(signalChannel, os.Interrupt)

	go func() {
		if err := cmd.Execute(ctx); err != nil {
			exitChannel <- 1
		}

		exitChannel <- 0
	}()

	select {
	case <-signalChannel:
		slog.Warn("Received interrupt signal")
		cancel()
	case code := <-exitChannel:
		os.Exit(code)
	}
}
