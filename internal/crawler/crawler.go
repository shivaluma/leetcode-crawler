package crawler

import (
	"time"

	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
	"leetcode-submission-crawler/config"
	"leetcode-submission-crawler/logger"
)

type Crawler struct {
	config    *config.Config
	scheduler *gocron.Scheduler
}

func NewCrawler(config *config.Config) *Crawler {
	return &Crawler{
		config:    config,
		scheduler: gocron.NewScheduler(time.UTC),
	}
}

func (c *Crawler) RegisterCron(cronExpresion string, job func()) {
	_, err := c.scheduler.Cron(cronExpresion).Do(job)
	if err != nil {
		logger.Log.Error("crawler.RegisterCron", zap.Error(err))
		return
	}
}
