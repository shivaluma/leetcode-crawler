package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/golang-migrate/migrate/v4"
	"go.uber.org/zap"
	"leetcode-submission-crawler/config"
	"leetcode-submission-crawler/database"
	"leetcode-submission-crawler/internal/crawler"
	"leetcode-submission-crawler/internal/github"
	"leetcode-submission-crawler/internal/infrastructure/repository"
	"leetcode-submission-crawler/internal/leetcode"
	"leetcode-submission-crawler/internal/workerpool"
	"leetcode-submission-crawler/logger"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

/*
There will be two way to trigger for an update:
  1. by cron expression: we will provide a cron expression to trigger the update job
  2. by rest api: we will provide a rest api to trigger the update job

  The crawler app will receive a channel, that contains the command to trigger the update job.
*/

func main() {
	logger.InitLogger()

	var cfg *config.Config

	cfg = config.ReadConfig()

	logger.Log.Info("starting migration")

	m, err := migrate.New(
		"file://./migrations",
		database.PostgresConnectionString(cfg.Postgres))

	if err != nil {

		logger.Log.Fatal("database.MigrationPrepare", zap.Error(err))

	}

	err = m.Up()
	if err != nil && err.Error() != "no change" {
		logger.Log.Fatal("database.Migration", zap.Error(err))
	}

	logger.Log.Info("starting scheduler")

	scheduler := gocron.NewScheduler(time.UTC)

	lc := leetcode.NewLeetcode(cfg)

	gh := github.NewGithub(cfg)

	wp := workerpool.New(1)

	pgx := database.NewDBClient(cfg)

	lutRepo := repository.NewLastUpdatedTimestampRepository(pgx)
	combinedRepo := repository.NewRepository(lutRepo)
	newCrawler := crawler.NewCrawler(cfg, wp, lc, gh, combinedRepo)

	scheduler.Every(1).Hour().Do(newCrawler.Crawl)

	scheduler.StartAsync()

	go wp.Start()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	for {
		select {
		case <-signalChan:
			logger.Sync()
			return
		}
	}

}
