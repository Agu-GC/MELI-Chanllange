package pkg

import (
	"os"

	"github.com/rs/zerolog"
)

type Logger interface {
	Info(msg string, fields map[string]any)
	Error(msg string, fields map[string]any)
	Fatal(msg string, fields map[string]any)
	With(fields map[string]any) Logger
}

type ZeroLogger struct {
	logger zerolog.Logger
}

func NewLogger(serviceName string) Logger {
	return &ZeroLogger{
		logger: zerolog.New(os.Stdout).
			With().
			Timestamp().
			Str("service", serviceName).
			Logger(),
	}
}

func (l *ZeroLogger) Info(msg string, fields map[string]any) {
	l.logger.Info().Caller(1).Fields(fields).Msg(msg)
}

func (l *ZeroLogger) Error(msg string, fields map[string]any) {
	l.logger.Error().Caller(1).Fields(fields).Msg(msg)
}

func (l *ZeroLogger) Fatal(msg string, fields map[string]any) {
	l.logger.Fatal().Caller(1).Fields(fields).Msg(msg)
}

func (l *ZeroLogger) With(fields map[string]any) Logger {
	return &ZeroLogger{
		logger: l.logger.With().Fields(fields).Logger(),
	}
}
