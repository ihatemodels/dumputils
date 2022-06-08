/*
 * Copyright (c) 2022. Dumputils Authors
 *
 * https://opensource.org/licenses/MIT
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

// Package log provides a global zerolog instance.
package log

import (
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

func Init(level string, logType string) {
	if logType == "human" {
		cWriter := &consoleWriter{
			level: getLogLevel(level),
			out:   os.Stdout,
		}
		logger = zerolog.New(cWriter).
			With().
			Timestamp().
			Logger()
	} else {
		zerolog.SetGlobalLevel(getLogLevel(level))
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
