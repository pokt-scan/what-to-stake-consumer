package wtsc

import "os"

func Init() {
	// define default logger to use before load config and override it
	Logger = GetDefaultLogger()

	version := os.Getenv("VERSION")
	if version == "" {
		version = "0.0.0"
	}
	Logger.Info().Str("version", version).Msg("initializing wtsc")

	cfg := LoadConfig()

	AppConfig = cfg
}
