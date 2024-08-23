package wtsc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pokt-scan/wtsc/generated"
	"github.com/robfig/cron/v3"
	"os"
	"path/filepath"
	"time"
)

func writeResults(result *generated.GetWhatToStakeResponse) {
	// Convert the struct to pretty-printed JSON
	prettyJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		Logger.Error().Err(err).Msg("error marshalling to JSON what to stake result")
		return
	}

	// Get the current date and time
	currentTime := time.Now()

	// Format the date and time to a string
	dateTimeString := currentTime.Format("20060102_150405")

	// Define the file name with the formatted date and time
	fileName := fmt.Sprintf("file_%s.json", dateTimeString)

	// Concatenate paths to form the full file path
	fullPath := filepath.Join(AppConfig.ResultsPath, fileName)

	// Write the JSON to a file
	file, err := os.Create(fullPath)
	if err != nil {
		Logger.Error().Err(err).Str("path", fullPath).Msg("error creating file")
		return
	}
	defer func(file *os.File) {
		e := file.Close()
		if e != nil {
			Logger.Error().Err(e).Str("path", fullPath).Msg("error closing file results file")
		}
	}(file)

	_, err = file.Write(prettyJSON)
	if err != nil {
		Logger.Error().Err(err).Str("path", fullPath).Msg("error writing to file")
		return
	}

	Logger.Info().Str("path", fullPath).Msg("writing results to file")
}

func evaluationJob() {
	Logger.Info().Msg("running evaluation")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(AppConfig.MaxTimeout)*time.Millisecond)
	defer cancel()

	input := generated.WtsProcessRequestInput{
		Domain:               AppConfig.Domain,
		Service_pool:         AppConfig.ServicePool,
		Min_increase_percent: AppConfig.MinIncreasePercent,
		Stake_weight:         int(AppConfig.StakeWeight),
		Min_service_stake:    AppConfig.MinServiceStake.CastToGqlType(),
		Time_period:          int(AppConfig.TimePeriod),
	}

	Logger.Info().Msg("calling what to stake service")
	resp, err := generated.GetWhatToStake(ctx, *POKTscanApiClient, input)

	if err != nil {
		Logger.Error().Err(err).Str("service", AppConfig.POKTscanApi).Msg("failed to call what to stake service")
		return
	}

	if IsEmptyString(AppConfig.ResultsPath) {
		resultStr, e := json.Marshal(resp)
		if e != nil {
			Logger.Error().Err(e).Msg("failed to marshal results")
		} else {
			Logger.Debug().Str("result", string(resultStr)).Msg("what to stake results")
		}
	} else {
		writeResults(resp)
	}

	if AppConfig.DryMode {
		Logger.Info().Msg("DRY MODE is on, omitting stake transactions.")
		return
	}

	if !resp.GetWhatToStake.Do_update {
		Logger.Info().Msg("What-To-Stake thinks you does not need to update yet.")
		return
	}

	// create a group inside worker pool because it allows just waiting without a stop
	group := WorkerPool.Group()

	for _, wtsServicer := range resp.GetWhatToStake.Servicers {
		if signer, ok := ServicersMap.Load(wtsServicer.Address); !ok {
			Logger.Warn().Str("address", wtsServicer.Address).Msg("Failed to find signer")
			continue
		} else {
			group.Submit(StakeServicer(signer, &wtsServicer))
		}
	}

	// Stop group pool and wait for all submitted tasks to complete
	group.Wait()
}

func Schedule(frequency string) (entry cron.EntryID, err error) {
	Logger.Info().Msg("preparing cron job")
	CronJob = cron.New()
	// Define the job
	entry, err = CronJob.AddFunc(frequency, evaluationJob)
	Logger.Debug().Int("schedule_id", int(entry)).Msg("scheduled job detail")
	// Start the cron job
	CronJob.Start()
	return
}

func ReSchedule(frequency string) error {
	if WorkerPool.WaitingTasks() > 0 {
		// prefer to delay the schedule update to the moment where there are not waiting tasks
		return errors.New("unable to update schedule due to waiting jobs")
	}

	CronJob.Stop()

	entryId, err := Schedule(frequency)

	if err != nil {
		return err
	}

	Logger.Debug().Int("NewScheduleID", int(entryId)).Msg("re-scheduled job")

	return nil
}
