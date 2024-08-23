package main

import (
	"github.com/pokt-scan/wtsc"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func Run() {
	defer func() {
		if r := recover(); r != nil {
			log.Trace().Stack().Timestamp().
				Interface("recover", r).
				Msg("Recovered from panic")
		}
	}()

	// hardcode time for now - could be replaced by env var?
	// Default duration
	defaultDuration := 30 * time.Second

	// Get the environment variable
	envReloadSeconds := os.Getenv("RELOAD_SECONDS")

	// Parse the environment variable
	reloadSeconds, err := strconv.Atoi(envReloadSeconds)
	var sleepDuration time.Duration

	if err != nil || reloadSeconds <= 0 {
		// Use the default duration if parsing fails or if the value is non-positive
		sleepDuration = defaultDuration
	} else {
		// Use the parsed duration
		sleepDuration = time.Duration(reloadSeconds) * time.Second
	}

	log.Info().Dur("duration", sleepDuration).Msg("config reload")

	for {
		time.Sleep(sleepDuration)
		wtsc.ReloadConfig()
	}
}

func main() {
	cfg := wtsc.LoadConfig()

	if valid, wrongKeys := wtsc.ValidateConfig(cfg); !valid {
		log.Fatal().
			Str("path", wtsc.GetConfigFilePath()).
			Int("count", len(wrongKeys)).
			Strs("keys", wrongKeys).
			Msg("loaded config.json contains errors")
	}

	wtsc.AppConfig = cfg

	// Configure logger
	wtsc.ConfigLogger(wtsc.AppConfig.LogLevel, wtsc.AppConfig.LogFormat)

	// Configure http client
	wtsc.NewHttpClient(wtsc.AppConfig.POKTscanApiToken, wtsc.AppConfig.MaxRetries, wtsc.AppConfig.MaxTimeout)

	// Create POKTscan Client
	wtsc.NewPOKTscanClient(wtsc.AppConfig.POKTscanApi)

	// Create PocketRpcProvider
	wtsc.NewPocketRpcProvider(wtsc.AppConfig.PocketRPC, wtsc.AppConfig.MaxRetries, wtsc.AppConfig.MaxTimeout)

	// Initialize the worker pool
	wtsc.NewWorker(cfg.MaxWorkers, uint(len(cfg.ServicerKeys)))

	// Initialize the servicers map
	wtsc.NewSignerMap(cfg.ServicerKeys)

	// Initialize the cron job
	_, err := wtsc.Schedule(cfg.Schedule)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to schedule cron job")
	}

	// Create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1)

	// Notify the channel on interrupt and termination signals
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Info().Msg("Starting What-To-Stake Consumer...")
	go Run()

	// Block until we receive a signal
	sig := <-sigChan
	log.Info().Str("signal", sig.String()).Msg("Received signal. Exiting...")
	log.Info().Msg("Shutting down...")
	// stop cron job schedule another one
	wtsc.CronJob.Stop()
	// wait for any in progress job.
	log.Debug().Uint64("waiting_tasks", wtsc.WorkerPool.WaitingTasks()).Msg("Shutting down Workers...")
	wtsc.WorkerPool.StopAndWait()
	log.Info().Msg("See you later, baby!")
	os.Exit(0)
}
