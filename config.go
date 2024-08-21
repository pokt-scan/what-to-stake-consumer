package main

import (
	"encoding/json"
	"github.com/alitto/pond"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
)

type Config struct {
	// Domain is Servicer domain
	Domain string `json:"domain"`
	// ChainPool list of the chains that are supported by the node runner.
	ChainPool []string `json:"chain_pool"`
	// ServicerKeys list of the private keys to use for sign nodes
	ServicerKeys []string `json:"servicer_keys"`
	// LogLevel is the level of logging
	LogLevel string `json:"log_level"`
	// Schedule is the cron schedule. Allow @every 1m or cron text.
	// Refer to: https://pkg.go.dev/github.com/robfig/cron#hdr-CRON_Expression_Format
	Schedule string `json:"schedule"`
	// Worker Pool config
	MaxWorkers uint64 `json:"max_workers"`
}

func NewWorker(maxWorkers, maxCapacity uint64) {
	workerPool = pond.New(int(maxWorkers), int(maxCapacity))
}

func ValidateConfig(cfg *Config) error {
	// TODO: add as many validation as it needs here
	return nil
}

func LoadConfig() *Config {
	cfg := Config{}
	var configPath = "./config.json"
	var jsonFile *os.File
	var bz []byte
	if _, err := os.Stat(configPath); err != nil && os.IsNotExist(err) {
		log.Fatal().Str("configPath", configPath).Msg("config file not found")
	}
	// reopen the file to read into the variable
	jsonFile, err := os.OpenFile(configPath, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal().Str("configPath", configPath).Err(err).Msg("failed to open config file")
	}
	bz, err = io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal().Str("configPath", configPath).Err(err).Msg("failed to read config file")
	}
	err = jsonFile.Close()
	if err != nil {
		log.Fatal().Str("configPath", configPath).Err(err).Msg("failed to close config file")
	}
	err = json.Unmarshal(bz, &cfg)
	if err != nil {
		log.Fatal().Str("configPath", configPath).Err(err).Msg("failed to unmarshal config file")
	}

	err = ValidateConfig(&cfg)
	if err != nil {
		log.Fatal().Str("configPath", configPath).Err(err).Msg("failed to validate config file")
	}

	return &cfg
}

func ReloadConfig(cfg *Config) {
	newCfg := LoadConfig()

	if cfg.Domain != newCfg.Domain {
		log.Debug().Msg("change detected at config.domain")
		cfg.Domain = newCfg.Domain
	}

	if diff := GetSliceDiff(cfg.ChainPool, newCfg.ChainPool); len(diff) > 0 {
		log.Debug().Msg("change detected at config.chain_pool")
		cfg.ChainPool = newCfg.ChainPool
	}

	if diff := GetSliceDiff(cfg.ServicerKeys, newCfg.ServicerKeys); len(diff) > 0 {
		log.Debug().Msg("change detected at config.servicer_keys")
		cfg.ServicerKeys = newCfg.ServicerKeys
		if len(cfg.ServicerKeys) < len(newCfg.ServicerKeys) {
		}
	}

	if cfg.LogLevel != newCfg.LogLevel {
		log.Debug().Msg("change detected at config.log_level")
		// here it will already be validated by LoadConfig but just in case.
		newLvl, err := zerolog.ParseLevel(cfg.LogLevel)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse log level")
		} else {
			cfg.LogLevel = newCfg.LogLevel
			zerolog.SetGlobalLevel(newLvl)
		}
	}

	if cfg.Schedule != newCfg.Schedule {
		log.Debug().Msg("change detected at config.schedule")
		err := ReSchedule(cfg.Schedule)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse log level")
		} else {
			// assign to cfg reference the new value.
			cfg.Schedule = newCfg.Schedule
		}
	}

	if cfg.MaxWorkers != newCfg.MaxWorkers {
		waitingTasks := workerPool.WaitingTasks()
		if waitingTasks > 0 {
			log.Error().Uint64("waiting_tasks", waitingTasks).Msg("unable to update worker pool min/max workers due to it has waiting tasks. will be retried on next round.")
		} else {
			workerPool.StopAndWait()
			// max capacity will always be the same of servicer keys
			NewWorker(newCfg.MaxWorkers, uint64(len(newCfg.ServicerKeys)))
		}
	}
}
