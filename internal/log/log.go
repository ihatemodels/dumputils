// Package log provides a global zerolog instance.
package log

import (
	"github.com/ihatemodels/dumputils/internal/config"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog"
	"os"
	"strings"
)

// logger is the global logger instance.
var logger = zerolog.New(os.Stdout).
	With().
	Timestamp().
	Logger()

func Configure() {
	if config.App.Log.Type == "human" {
		cWriter := &consoleWriter{
			level: getLogLevel(config.App.Log.Level),
			out:   os.Stdout,
		}
		logger = zerolog.New(cWriter).
			With().
			Timestamp().
			Logger()
	} else {
		zerolog.SetGlobalLevel(getLogLevel(config.App.Log.Level))
	}
}

// GetLogLevel returns zerolog.Level by given string
func getLogLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}

func Debugf(format string, args ...any) {
	logger.Debug().Msgf(format, args...)
}

func Infof(format string, args ...any) {
	logger.Info().Msgf(format, args...)
}

func Warnf(format string, args ...any) {
	logger.Info().Msgf(format, args...)
}

func Errorf(err error, format string, args ...any) {
	logger.Error().
		Err(err).
		Str("stacktrace", eris.ToString(err, true)).
		Msgf(format, args...)
}
