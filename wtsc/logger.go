package wtsc

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"os"
	"time"
)

const (
	LogJsonFormat = "json"
	LogTextFormat = "text"
)

// ZerologLeveledLogger is an implementation of retryablehttp.LeveledLogger using zerolog
type ZerologLeveledLogger struct {
	logger zerolog.Logger
}

// NewZerologLeveledLogger creates a new ZerologLeveledLogger
func NewZerologLeveledLogger(l zerolog.Logger) *ZerologLeveledLogger {
	return &ZerologLeveledLogger{
		logger: l,
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

func GetDefaultLogger() zerolog.Logger {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	// prevent slow down the process due to log write on console
	wr := diode.NewWriter(os.Stdout, 1000, 10*time.Millisecond, func(missed int) {
		log.Warn().Msgf("Logger Dropped %d messages", missed)
	})

	return zerolog.New(wr).With().
		Timestamp(). // add timestamp that will be unix format (faster)
		Caller().    // add caller to logs
		Stack().     // add stack trace to errors only
		Logger().
		// prevent crazy amount of logs on debug
		Sample(zerolog.LevelSampler{
			DebugSampler: &zerolog.BurstSampler{
				Burst:       5,
				Period:      1 * time.Second,
				NextSampler: &zerolog.BasicSampler{N: 100},
			},
		}).
		Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

func ConfigLogger(level, format string) {
	// level is already parse on ValidateConfig
	newLvl, _ := zerolog.ParseLevel(level)

	Logger = GetDefaultLogger().Level(newLvl)

	// Override the output to console
	if format == LogTextFormat {
		// replace global logger format
		Logger = Logger.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		Logger.Info().Str("format", format).Msg("switch logger format")
	}
}
