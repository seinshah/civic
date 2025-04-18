package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/seinshah/civic/internal/command"
	"github.com/seinshah/civic/internal/pkg/logger"
	"github.com/seinshah/civic/internal/pkg/version"
)

// Version will be changed at build time using ldflags.
// The new value will be determined in the CI pipeline
// using the git tag.
//
//nolint:gochecknoglobals
var Version = "0.1.0-dev"

func main() {
	logger.SetUp()

	cmd := command.NewCommand(Version)

	ctx, ctxCanceler := context.WithCancel(context.Background())

	signalChannel := make(chan os.Signal, 1)
	exitChannel := make(chan int, 1)

	signal.Notify(signalChannel, os.Interrupt)

	checkVersions(ctx)

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

func checkVersions(ctx context.Context) {
	var (
		latestVersion  *version.Semantic
		currentVersion *version.Semantic
		errLV          error
		errCV          error
	)

	latestVersion, errLV = version.ParseFromGithub(ctx)

	if Version == "latest" && errLV == nil {
		Version = latestVersion.String()
		currentVersion = latestVersion
	} else if Version != "latest" {
		currentVersion, errCV = version.Parse(Version)
		if errCV != nil {
			slog.Error("current app version is not valid an might lead to unexpected behavior", "error", errCV)
		}
	}

	if latestVersion == nil || currentVersion == nil {
		slog.Warn("Newer version is available.", "latest", latestVersion, "current", currentVersion)

		return
	}

	if currentVersion.Equal(latestVersion) {
		slog.Info("you are using the latest version", "version", currentVersion)
	}

	slog.Warn("you are using an outdated version", "latest", latestVersion, "current", currentVersion)
}
