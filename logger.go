package wtsc

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

// ZerologLeveledLogger is an implementation of retryablehttp.LeveledLogger using zerolog
type ZerologLeveledLogger struct {
	logger zerolog.Logger
}

// NewZerologLeveledLogger creates a new ZerologLeveledLogger
func NewZerologLeveledLogger() *ZerologLeveledLogger {
	return &ZerologLeveledLogger{
		logger: log.Output(zerolog.ConsoleWriter{Out: os.Stderr}),
	}
}

// Error logs a message at level Error.
func (l *ZerologLeveledLogger) Error(msg string, keysAndValues ...interface{}) {
	l.logger.Error().Fields(keysAndValues).Msg(msg)
}

// Info logs a message at level Info.
func (l *ZerologLeveledLogger) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Info().Fields(keysAndValues).Msg(msg)
}

// Debug logs a message at level Debug.
func (l *ZerologLeveledLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.logger.Debug().Fields(keysAndValues).Msg(msg)
}

// Warn logs a message at level Warn.
func (l *ZerologLeveledLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.logger.Warn().Fields(keysAndValues).Msg(msg)
}

func ConfigLogger(level string) {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	newLvl, err := zerolog.ParseLevel(level)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse log level")
	} else {
		zerolog.SetGlobalLevel(newLvl)
	}
}
