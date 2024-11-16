package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

type options struct {
	level   slog.Leveler
	noColor bool
}

type Option func(*options)

var tintOptions *tint.Options

func WithLevel(level slog.Level) Option {
	return func(o *options) {
		o.level = level
	}
}

func WithNoColor() Option {
	return func(o *options) {
		o.noColor = true
	}
}

func SetUp(opts ...Option) {
	setTintOptions()

	if len(opts) == 0 {
		setDefaultLogger()

		return
	}

	newOpts := options{
		level: tintOptions.Level,
	}

	for _, opt := range opts {
		opt(&newOpts)
	}

	if newOpts.noColor {
		tintOptions.NoColor = true
	}

	tintOptions.Level = newOpts.level

	setDefaultLogger()
}

func setTintOptions() {
	if tintOptions == nil {
		tintOptions = &tint.Options{
			Level: slog.LevelInfo,
		}
	}
}

func setDefaultLogger() {
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, tintOptions),
	))
}
