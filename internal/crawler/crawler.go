package crawler

import (
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"
	"leetcode-submission-crawler/config"
	"leetcode-submission-crawler/internal/github"
	"leetcode-submission-crawler/internal/infrastructure/repository"
	"leetcode-submission-crawler/internal/leetcode"
	"leetcode-submission-crawler/internal/workerpool"
	"leetcode-submission-crawler/logger"
)

type Crawler interface {
	Crawl()
}

type Impl struct {
	wp     *workerpool.WorkerPool
	lc     leetcode.UseCase
	gh     github.UseCase
	config *config.Config
	repo   *repository.Repository
}

func (i Impl) Crawl() {

	shouldContinue := true
	lut, err := i.repo.Lastupdatedtimestamp.Get(i.config.Leetcode.Username)

	logger.Log.Info("repo.Lastupdatedtimestamp.Get", zap.Any("lut", lut))

	if err != nil {
		logger.Log.Error("repo.Lastupdatedtimestamp.Get", zap.Error(err))
	}

	recents := i.lc.GetRecentSubmissions(20)

	if lut != nil {
		if len(recents) > 0 {
			// last submissions
			tsInt, _ := strconv.Atoi(recents[0].Timestamp)
			if time.Unix(int64(tsInt), 0).Before(lut.LastUpdatedAt) {
				shouldContinue = false
			}
		}

		if !shouldContinue {
			logger.Log.Info("crawler.Crawl not found any new submission. Skipping....", zap.String("username", i.config.Leetcode.Username))
			return
		}
	}

	for idx := len(recents) - 1; idx >= 0; idx-- {

		logger.Log.Info("current idx", zap.Int("idx", idx))

		task := func() {
			currentIdx := idx
			recent := recents[currentIdx]

			if lut != nil {
				tsInt, _ := strconv.ParseInt(recent.Timestamp, 10, 64)
				ts := time.Unix(int64(tsInt), 0)

				if ts.Before(lut.LastUpdatedAt) {
					logger.Log.Info("This submission might be updated, skip...", zap.Any("recent", recent))
					return
				}
			}

			submissionId, _ := strconv.ParseInt(recent.ID, 10, 64)

			submission := i.lc.GetSubmissionDetail(submissionId)

			if submission == nil {
				logger.Log.Error("crawler.Crawl", zap.String("detail", "cannot get submission detail"), zap.String("username", i.config.Leetcode.Username), zap.String("id", recent.ID))
				return
			}

			if submission.Question == nil {
				logger.Log.Error("crawler.Crawl", zap.String("detail", "cannot get question detail, need update token"))
				return
			}

			questionContent := i.lc.GetQuestionContent(recent.TitleSlug)

			if questionContent == nil {
				logger.Log.Error("crawler.Crawl", zap.String("detail", "cannot get question content"), zap.String("username", i.config.Leetcode.Username), zap.String("id", recent.ID))
				return
			}

			// get SHA of readme

			if err != nil {
				logger.Log.Error("crawler.GetSHA", zap.Error(err))
			}

			// create new git commit and push file
			folderName := submission.Question.QuestionId + "-" + recent.TitleSlug
			fileName := folderName + "." + leetcode.ExtensionNameMapping(recent.Lang)

			msg := fmt.Sprintf("Time: %s (%s) - Space: %s (%s) - Leetcrawl", recent.Runtime, fmt.Sprintf("%.0f%%", submission.RuntimePercentile), recent.Memory, fmt.Sprintf("%.0f%%", submission.MemoryPercentile))

			// create readme
			readmeSHA, err := i.gh.GetSHA(i.config.Github.Username, i.config.Github.Repo, folderName+"/README.md")
			logger.Log.Info("crawler.GetSHA", zap.Any("readmeSHA", readmeSHA))
			err = i.gh.CreateNewCommit(i.config.Github.Username, i.config.Github.Repo, folderName, "README.md", "Attach README - LeetCrawl", questionContent.Content, readmeSHA)
			if err != nil {
				logger.Log.Error("crawler.CreateNewCommit", zap.Error(err))
				return
			}

			// get SHA of code file.
			codeSHA, _ := i.gh.GetSHA(i.config.Github.Username, i.config.Github.Repo, folderName+"/"+fileName)
			err = i.gh.CreateNewCommit(i.config.Github.Username, i.config.Github.Repo, folderName, fileName, msg, submission.Code, codeSHA)
			if err != nil {
				logger.Log.Error("crawler.CreateNewCommit", zap.Error(err))
				return
			}

			// update last updated timestamp
			i.repo.Lastupdatedtimestamp.Upsert(i.config.Leetcode.Username, i.config.Github.Username, time.Unix(int64(submission.Timestamp), 0))
		}

		i.wp.Enqueue(task)
	}

}

func NewCrawler(config *config.Config, wp *workerpool.WorkerPool, lc leetcode.UseCase, gh github.UseCase, repo *repository.Repository) Crawler {

	return &Impl{
		wp:     wp,
		lc:     lc,
		gh:     gh,
		config: config,
		repo:   repo,
	}
}
