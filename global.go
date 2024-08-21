package main

import (
	"github.com/alitto/pond"
	"github.com/robfig/cron/v3"
)

var (
	workerPool *pond.WorkerPool
	cronJob    *cron.Cron
	config     *Config
)
