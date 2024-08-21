package main

import (
	"errors"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

func evaluationJob() {
	// create a group inside worker pool because it allows just waiting without a stop
	group := workerPool.Group()
	// 1. call what-to-stake
	// 2. evaluate response
	// 3. add worker pool jobs for each address to run the stake
	group.Submit(func() {
		// stake
	})
	// Stop group pool and wait for all submitted tasks to complete
	group.Wait()
}

func Schedule(frequency string) (cron.EntryID, error) {
	cronJob = cron.New()
	return cronJob.AddFunc(frequency, evaluationJob)
}

func ReSchedule(frequency string) error {
	workerPool.WaitingTasks()
	if workerPool.WaitingTasks() > 0 {
		// prefer to delay the schedule update to the moment where there are not waiting tasks
		return errors.New("unable to update schedule due to waiting jobs")
	}

	cronJob.Stop()

	entryId, err := Schedule(frequency)

	if err != nil {
		return err
	}

	log.Debug().Int("ScheduleID", int(entryId)).Msg("re-scheduled job")

	return nil
}
