package wtsc

import (
	"errors"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

func evaluationJob() {
	log.Info().Msg("Running evaluation job")

	// @TODO: Uncomment once WTS exists
	//resp, err := HttpRetryableClient.Post(AppConfig.WhatToStakeService, "application/json", WTSParams{
	//	Domain:             AppConfig.Domain,
	//	ServicePool:        AppConfig.ChainPool,
	//	MinIncreasePercent: AppConfig.MinIncreasePercent,
	//	StakeWeight:        AppConfig.StakeWeight,
	//	MinServiceStake:    AppConfig.MinServiceStake,
	//})
	//if err != nil {
	//	log.Error().Err(err).Str("service", AppConfig.WhatToStakeService).Msg("Failed to call what to stake service")
	//}
	//if resp.StatusCode != http.StatusOK {
	//	log.Error().Err(errors.New(resp.Status)).Int("code", resp.StatusCode).Str("service", AppConfig.WhatToStakeService).Msg("Failed to call what to stake service")
	//}
	// @TODO: Update WTSResponse struct with the right fields and update code below
	whatToStakeResponse := WTSResponse{
		DoUpdate: true,
		Servicers: []*WTSService{
			{
				Address: "e5e87a03577606b9d8be9456c0424f228b06893d",
				Chains:  []string{"0021"},
			},
		},
	}

	if !whatToStakeResponse.DoUpdate {
		log.Info().Msg("What-To-Stake thinks you does not need to update yet.")
		return
	}

	// create a group inside worker pool because it allows just waiting without a stop
	group := WorkerPool.Group()

	for _, wtsServicer := range whatToStakeResponse.Servicers {
		if signer, ok := ServicersMap.Load(wtsServicer.Address); !ok {
			log.Warn().Str("address", wtsServicer.Address).Msg("Failed to find signer")
			continue
		} else {
			group.Submit(StakeServicer(signer, wtsServicer))
		}
	}

	// Stop group pool and wait for all submitted tasks to complete
	group.Wait()
}

func Schedule(frequency string) (entry cron.EntryID, err error) {
	CronJob = cron.New()
	// Define the job
	entry, err = CronJob.AddFunc(frequency, evaluationJob)
	log.Debug().Int("ScheduleID", int(entry)).Msg("scheduled job")
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

	log.Debug().Int("NewScheduleID", int(entryId)).Msg("re-scheduled job")

	return nil
}
