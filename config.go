package wtsc

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
)

// ValidateConfig validates the provided configuration struct.
// TODO: add as many validation as it needs here
func ValidateConfig(cfg *Config) error {
	// TODO: add as many validation as it needs here
	return nil
}

// LoadConfig loads and returns the configuration by reading the `configPath` file.
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

// ReloadConfig updates the provided configuration based on changes detected in a new configuration.
func ReloadConfig() {
	log.Debug().Msg("Reloading configuration")
	newCfg := LoadConfig()

	if AppConfig.Domain != newCfg.Domain {
		log.Debug().Msg("change detected at config.domain")
		AppConfig.Domain = newCfg.Domain
	}

	if diff := GetStrSliceDiff(AppConfig.ChainPool, newCfg.ChainPool); len(diff) > 0 {
		log.Debug().Msg("change detected at config.chain_pool")
		AppConfig.ChainPool = newCfg.ChainPool
	}

	if diff := GetStrSliceDiff(AppConfig.ServicerKeys, newCfg.ServicerKeys); len(diff) > 0 {
		log.Debug().Msg("change detected at config.servicer_keys")
		AppConfig.ServicerKeys = newCfg.ServicerKeys
		UpdateServicers(AppConfig.ServicerKeys)
	}

	if AppConfig.StakeWeight != newCfg.StakeWeight {
		log.Debug().Msg("change detected at config.stake_weight")
		AppConfig.StakeWeight = newCfg.StakeWeight
	}

	if AppConfig.MinIncreasePercent != newCfg.MinIncreasePercent {
		log.Debug().Msg("change detected at config.min_increase_percent")
		AppConfig.MinIncreasePercent = newCfg.MinIncreasePercent
	}

	if AppConfig.StakeWeight != newCfg.StakeWeight {
		log.Debug().Msg("change detected at config.stake_weight")
		AppConfig.StakeWeight = newCfg.StakeWeight
	}

	if updated, removed, added := GetServiceStakeSliceDiff(AppConfig.MinServiceStake, newCfg.MinServiceStake); len(updated) > 0 || len(removed) > 0 || len(added) > 0 {
		log.Debug().Msg("change detected at config.min_service_stake")
		AppConfig.MinServiceStake = newCfg.MinServiceStake
	}

	if AppConfig.LogLevel != newCfg.LogLevel {
		log.Debug().Msg("change detected at config.log_level")
		// here it will already be validated by LoadConfig but just in case.
		newLvl, err := zerolog.ParseLevel(AppConfig.LogLevel)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse log level")
		} else {
			AppConfig.LogLevel = newCfg.LogLevel
			zerolog.SetGlobalLevel(newLvl)
		}
	}

	if AppConfig.Schedule != newCfg.Schedule {
		log.Debug().Msg("change detected at config.schedule")
		err := ReSchedule(AppConfig.Schedule)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse log level")
		} else {
			// assign to AppConfig reference the new value.
			AppConfig.Schedule = newCfg.Schedule
		}
	}

	if AppConfig.MaxWorkers != newCfg.MaxWorkers {
		waitingTasks := WorkerPool.WaitingTasks()
		if waitingTasks > 0 {
			log.Error().Uint64("waiting_tasks", waitingTasks).Msg("unable to update worker pool min/max workers due to it has waiting tasks. will be retried on next round.")
		} else {
			WorkerPool.StopAndWait()
			// max capacity will always be the same of servicer keys
			NewWorker(newCfg.MaxWorkers, uint64(len(newCfg.ServicerKeys)))
		}
	}

}
