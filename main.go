package main

import (
	"github.com/rs/zerolog"
)

func main() {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	cfg := LoadConfig()
	config = &cfg
}
