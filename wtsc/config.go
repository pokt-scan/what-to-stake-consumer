package wtsc

import (
	"encoding/json"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"io"
	"os"
	"path/filepath"
)

type UpdateKeys struct {
	s []string
}

var (
	ProjectRoot    string
	ConfigFilePath string
)

func (up *UpdateKeys) Add(key string) {
	if up.s == nil {
		up.s = make([]string, 0)
	}
	up.s = append(up.s, key)
}

func (up *UpdateKeys) Size() int {
	return len(up.s)
}

func (up *UpdateKeys) Values() []string {
	return up.s
}

// GetConfigFilePath returns the file path of the config.json file in the current working directory.
// If there is an error while getting the working directory, it returns an empty string.
// The file path is created by joining the working directory path and the file name "config.json".
func GetConfigFilePath() string {
	// Check if the environment variable is set
	ProjectRoot = os.Getenv("PROJECT_ROOT")
	if ProjectRoot == "" {
		// Fallback to current working directory if the env variable is not set
		var err error
		ProjectRoot, err = os.Getwd()
		if err != nil {
			Logger.Fatal().Err(err).Msg("failed to get working dir")
		}
	}

	fileName := os.Getenv("CONFIG_FILE")
	if IsEmptyString(fileName) {
		fileName = "config.json"
	}

	// Create a path to the file in the working directory
	ConfigFilePath = filepath.Join(ProjectRoot, fileName)

	return ConfigFilePath
}

// ValidateConfig validates the provided configuration struct.
func ValidateConfig(cfg *Config) (valid bool, errors []string) {
	Logger.Info().Msg("validating config file")

	if !IsValidHttpURI(cfg.POKTscanApi) {
		errors = append(errors, "poktscan_api")
	}

	if IsEmptyString(cfg.POKTscanApiToken) {
		errors = append(errors, "poktscan_api_token")
	}

	if cfg.NetworkID != "mainnet" && cfg.NetworkID != "testnet" {
		errors = append(errors, "network_id")
	}

	if fee, err := cfg.TxFee.Int64(); err != nil || fee <= 0 {
		errors = append(errors, "tx_fee")
	}

	if !IsValidDomain(cfg.Domain) {
		errors = append(errors, "domain")
	}

	if !IsValidChainPool(cfg.ServicePool) {
		errors = append(errors, "service_pool")
	}

	if !IsValidServicerList(cfg.ServicerKeys, cfg.DryMode) {
		errors = append(errors, "servicer_keys")
	}

	if cfg.StakeWeight < 1 || cfg.StakeWeight > 4 {
		errors = append(errors, "stake_weight")
	}

	if cfg.MinIncreasePercent <= 1 && cfg.MinIncreasePercent > 100 {
		errors = append(errors, "min_increase_percent")
	}

	if !IsValidMinServiceStake(cfg.MinServiceStake) {
		errors = append(errors, "min_service_stake")
	}

	if cfg.TimePeriod < 6 || cfg.TimePeriod > 48 {
		errors = append(errors, "time_period")
	}

	if !IsEmptyString(cfg.ResultsPath) && !IsWritableDirectory(filepath.Join(ProjectRoot, cfg.ResultsPath)) {
		// empty string disable the feature so it's ok been empty
		errors = append(errors, "results_path")
	}

	if _, err := zerolog.ParseLevel(cfg.LogLevel); err != nil {
		errors = append(errors, "log_level")
	}

	if cfg.LogFormat != LogTextFormat && cfg.LogFormat != LogJsonFormat {
		errors = append(errors, "log_format")
	}

	if _, err := cron.ParseStandard(cfg.Schedule); err != nil {
		errors = append(errors, "schedule")
	}

	if cfg.MaxWorkers <= 0 {
		errors = append(errors, "max_workers")
	}

	if !IsValidHttpURI(cfg.PocketRPC) {
		errors = append(errors, "pocket_rpc")
	}

	if cfg.MaxRetries < 0 {
		errors = append(errors, "max_retries")
	}

	// less than 1s is not allowed
	if cfg.MaxTimeout < 1000 {
		errors = append(errors, "max_retries")
	}

	valid = len(errors) == 0

	return
}

// LoadConfig loads and returns the configuration by reading the `configPath` file.
func LoadConfig() *Config {
	cfg := Config{}
	var configPath = GetConfigFilePath()
	Logger.Info().Str("path", configPath).Msg("reading config file")
	var jsonFile *os.File
	var bz []byte
	if _, err := os.Stat(configPath); err != nil && os.IsNotExist(err) {
		Logger.Fatal().Str("path", configPath).Msg("config file not found")
	}
	// reopen the file to read into the variable
	jsonFile, err := os.OpenFile(configPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		Logger.Fatal().Str("path", configPath).Err(err).Msg("failed to open config file")
	}
	bz, err = io.ReadAll(jsonFile)
	if err != nil {
		Logger.Fatal().Str("path", configPath).Err(err).Msg("failed to read config file")
	}
	err = jsonFile.Close()
	if err != nil {
		Logger.Fatal().Str("path", configPath).Err(err).Msg("failed to close config file")
	}
	err = json.Unmarshal(bz, &cfg)
	if err != nil {
		Logger.Fatal().Str("path", configPath).Err(err).Msg("failed to unmarshal config file")
	}

	if valid, wrongKeys := ValidateConfig(&cfg); !valid {
		Logger.Fatal().
			Str("path", GetConfigFilePath()).
			Int("count", len(wrongKeys)).
			Strs("keys", wrongKeys).
			Msg("loaded config.json contains errors")
	}

	return &cfg
}

// ReloadConfig updates the provided configuration based on changes detected in a new configuration.
func ReloadConfig() {
	Logger.Info().Msg("looking for changes on config file")

	newCfg := LoadConfig()

	var updatePOKTscanClient bool
	var updatePocketProvider bool
	var updateHttpClient bool
	var updateWorker bool
	var updateLogger bool
	var updateSigners bool
	var updateSchedule bool

	uk := UpdateKeys{}

	if AppConfig.DryMode != newCfg.DryMode {
		uk.Add("dry_mode")
		AppConfig.DryMode = newCfg.DryMode
	}

	if AppConfig.POKTscanApi != newCfg.POKTscanApi {
		uk.Add("poktscan_api")
		updatePOKTscanClient = true
	}

	if AppConfig.POKTscanApiToken != newCfg.POKTscanApiToken {
		uk.Add("poktscan_api_token")
		updatePOKTscanClient = true
		updateHttpClient = true
	}

	if AppConfig.NetworkID != newCfg.NetworkID {
		uk.Add("network_id")
		AppConfig.NetworkID = newCfg.NetworkID
	}

	if AppConfig.TxMemo != newCfg.TxMemo {
		uk.Add("tx_memo")
		AppConfig.TxMemo = newCfg.TxMemo
	}

	if AppConfig.TxFee != newCfg.TxFee {
		uk.Add("tx_fee")
		AppConfig.TxFee = newCfg.TxFee
	}

	if AppConfig.Domain != newCfg.Domain {
		uk.Add("domain")
		AppConfig.Domain = newCfg.Domain
	}

	if diff := GetStrSliceDiff(AppConfig.ServicePool, newCfg.ServicePool); len(diff) > 0 {
		uk.Add("service_pool")
		AppConfig.ServicePool = newCfg.ServicePool
	}

	if diff := GetStrSliceDiff(AppConfig.ServicerKeys, newCfg.ServicerKeys); len(diff) > 0 {
		uk.Add("servicer_keys")
		updateSigners = true
	}

	if AppConfig.StakeWeight != newCfg.StakeWeight {
		uk.Add("stake_weight")
		AppConfig.StakeWeight = newCfg.StakeWeight
	}

	if AppConfig.MinIncreasePercent != newCfg.MinIncreasePercent {
		uk.Add("min_increase_percent")
		AppConfig.MinIncreasePercent = newCfg.MinIncreasePercent
	}

	if updated, removed, added := GetServiceStakeSliceDiff(AppConfig.MinServiceStake, newCfg.MinServiceStake); len(updated) > 0 || len(removed) > 0 || len(added) > 0 {
		uk.Add("min_service_stake")
		AppConfig.MinServiceStake = newCfg.MinServiceStake
	}

	if AppConfig.TimePeriod != newCfg.TimePeriod {
		uk.Add("time_period")
		AppConfig.TimePeriod = newCfg.TimePeriod
	}

	if AppConfig.ResultsPath != newCfg.ResultsPath {
		uk.Add("results_path")
		AppConfig.ResultsPath = newCfg.ResultsPath
	}

	if AppConfig.LogLevel != newCfg.LogLevel {
		uk.Add("log_level")
		updateLogger = true
	}

	if AppConfig.LogFormat != newCfg.LogFormat {
		uk.Add("log_format")
		updateLogger = true
	}

	if AppConfig.Schedule != newCfg.Schedule {
		uk.Add("schedule")
		updateSchedule = true
	}

	if AppConfig.MaxWorkers != newCfg.MaxWorkers {
		uk.Add("max_workers")
		updateWorker = true
	}

	if AppConfig.PocketRPC != newCfg.PocketRPC {
		uk.Add("pocket_rpc")
		updatePocketProvider = true
	}

	if AppConfig.MaxRetries != newCfg.MaxRetries {
		uk.Add("max_retries")
		updateHttpClient = true
	}

	if AppConfig.MaxTimeout != newCfg.MaxTimeout {
		uk.Add("max_timeout")
		updateHttpClient = true
	}

	if uk.Size() == 0 {
		Logger.Debug().Msg("config file look the same as before.")
		return
	}

	Logger.Info().Strs("changed_keys", uk.Values()).Msg("changes are detected on config file, proceeding to update elements")

	if updateSigners {
		if index, err := UpdateServicers(AppConfig.ServicerKeys); err != nil {
			// update on the fly, no other config modify it
			Logger.Error().Err(err).Int("index", index).Msg("error updating signers")
		} else {
			// update all the props that could trigger it
			AppConfig.ServicerKeys = newCfg.ServicerKeys
		}
	}

	if updateWorker {
		Logger.Info().Msg("updating worker pool")
		waitingTasks := WorkerPool.WaitingTasks()
		if waitingTasks > 0 {
			Logger.Error().Uint64("waiting_tasks", waitingTasks).Msg("unable to update worker pool due to it has waiting tasks. will be retried on next round.")
		} else {
			WorkerPool.StopAndWait()
			// max capacity will always be the same of servicer keys
			NewWorker(newCfg.MaxWorkers, uint(len(newCfg.ServicerKeys)))
			// update all the props that could trigger it
			AppConfig.MaxWorkers = newCfg.MaxWorkers
		}
	}

	if updateSchedule {
		Logger.Info().Msg("updating schedule")
		err := ReSchedule(newCfg.Schedule) // update on the fly, no other config modify it
		if err != nil {
			Logger.Error().Err(err).Msg("failed to reschedule. check your schedule config.")
		} else {
			// assign to AppConfig reference the new value.
			AppConfig.Schedule = newCfg.Schedule
		}
	}

	if updateLogger {
		Logger.Info().Msg("updating logger")
		ConfigLogger(newCfg.LogLevel, newCfg.LogFormat)
		AppConfig.LogFormat = newCfg.LogFormat
		AppConfig.LogLevel = newCfg.LogLevel
	}

	if updateHttpClient {
		Logger.Info().Msg("updating http client")
		NewHttpClient(newCfg.POKTscanApiToken, newCfg.MaxRetries, newCfg.MaxTimeout)
		AppConfig.POKTscanApiToken = newCfg.POKTscanApiToken
		AppConfig.MaxRetries = newCfg.MaxRetries
		AppConfig.MaxTimeout = newCfg.MaxTimeout
	}

	if updatePocketProvider {
		Logger.Info().Msg("updating pocket rpc")
		NewPocketRpcProvider(newCfg.PocketRPC, newCfg.MaxRetries, newCfg.MaxTimeout)
		// update rpc url, but if retries or timeout was modified will be already update by previous if
		AppConfig.PocketRPC = newCfg.PocketRPC
	}

	if updatePOKTscanClient {
		Logger.Info().Msg("updating poktscan api")
		// this one also update the basic client.
		NewPOKTscanClient(newCfg.POKTscanApi)
		// update poktscan api url
		AppConfig.POKTscanApi = newCfg.POKTscanApi
	}
}
