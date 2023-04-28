package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"leetcode-submission-crawler/internal/entity"
	"leetcode-submission-crawler/logger"
)

type SubmissionRepositoryInterface interface {
	Get(leetcodeUsername string, submissionId string) (*entity.CrawlSubmission, error)
	Upsert(leetcodeUsername string, githubUsername string, problemId string, sha *string, lastUpdatedAt time.Time) error
}

type SubmissionRepository struct {
	db *pgxpool.Pool
}

func (l SubmissionRepository) Upsert(leetcodeUsername string, githubUsername string, problemId string, sha *string, lastUpdatedAt time.Time) error {

	query := `
		INSERT INTO submissions (leetcode_username, github_username, problem_id, submission_id, sha, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (leetcode_username, problem_id) DO UPDATE SET
			submission_id = $4,
			sha = $5,
			last_updated_at = $6
	`

	// Execute the upsert command with the parameters
	_, err := l.db.Exec(context.Background(), query, leetcodeUsername, githubUsername, problemId, sha, lastUpdatedAt)

	if err != nil {
		logger.Log.Error("SubmissionRepository.Upsert", zap.Error(err))
		return err
	}

	return nil

}

func (l SubmissionRepository) Get(leetcodeUsername string, submissionId string) (*entity.CrawlSubmission, error) {
	var cs entity.CrawlSubmission
	err := l.db.QueryRow(context.Background(), "select * from submissions where leetcode_username = $1 and submission_id = $2", leetcodeUsername, submissionId).Scan(&cs.Id, &cs.LeetcodeUsername, &cs.GithubUsername, &cs.ProblemId, &cs.SubmissionId, &cs.Sha, &cs.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		logger.Log.Error("SubmissionRepository.Get", zap.Error(err))
		return nil, err
	}

	return &cs, nil
}

func NewSubmissionRepository(db *pgxpool.Pool) SubmissionRepositoryInterface {
	return &SubmissionRepository{
		db: db,
	}
}
