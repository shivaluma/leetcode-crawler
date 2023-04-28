package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
	"leetcode-submission-crawler/config"
	"leetcode-submission-crawler/logger"
)

type Scheduler struct {
	config    *config.Config
	scheduler *gocron.Scheduler
}

func NewScheduler(config *config.Config) *Scheduler {
	return &Scheduler{
		config:    config,
		scheduler: gocron.NewScheduler(time.UTC),
	}
}

func (c *Scheduler) RegisterCron(cronExpression string, job func()) {
	_, err := c.scheduler.Cron(cronExpression).Do(job)
	if err != nil {
		logger.Log.Error("scheduler.RegisterCron", zap.Error(err))
		return
	}
}

func (c *Scheduler) StartAsync(job func()) {
	_, err := c.scheduler.Hour().Do(job)

	if err != nil {
		logger.Log.Error("scheduler.RegisterCron", zap.Error(err))
		return
	}
}

func (c *Scheduler) Start() {
	c.scheduler.StartAsync()
}
