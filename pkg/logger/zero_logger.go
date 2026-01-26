package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type ZeroLogger struct {
	log zerolog.Logger
}

func (l *ZeroLogger) Info(msg string, kv ...any)  { l.log.Info().Fields(kv).Msg(msg) }
func (l *ZeroLogger) Error(msg string, kv ...any) { l.log.Error().Fields(kv).Msg(msg) }
func (l *ZeroLogger) Debug(msg string, kv ...any) { l.log.Debug().Fields(kv).Msg(msg) }

func NewZeroLogger() *ZeroLogger {
	zl := zerolog.New(os.Stdout).With().Timestamp().CallerWithSkipFrameCount(4).Logger()
	return &ZeroLogger{
		log: zl,
	}
}
