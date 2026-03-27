package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

func New(isDebug bool) zerolog.Logger {
	logLevel := zerolog.InfoLevel
	if isDebug {
		logLevel = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(logLevel)

	var logger zerolog.Logger

	if isDebug {
		// Console writer (colored, pretty logs for development)
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
		consoleWriter.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		}
		consoleWriter.FormatCaller = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}
		consoleWriter.FormatMessage = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}

		logger = zerolog.New(consoleWriter).With().Timestamp().Caller().Logger()
	} else {
		// JSON logs for production
		logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
	}

	return logger
}
