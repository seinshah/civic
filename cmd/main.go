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

	ctx, ctxCanceler := context.WithCancel(context.Background())

	signalChannel := make(chan os.Signal, 1)
	exitChannel := make(chan int, 1)

	signal.Notify(signalChannel, os.Interrupt)

	go func() {
		slog.Debug("executing the command...")

		if err := cmd.Execute(ctx); err != nil {
			exitChannel <- 1
		}

		exitChannel <- 0
	}()

	var exitCode int

	select {
	case <-signalChannel:
		slog.Warn("Received interrupt signal")

		exitCode = 1
	case exitCode = <-exitChannel:
	}

	ctxCanceler()

	os.Exit(exitCode)
}
