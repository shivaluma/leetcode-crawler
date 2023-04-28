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

type LastUpdatedTimestampRepositoryInterface interface {
	Get(leetcodeId string) (*entity.LastUpdatedTimestamp, error)
	Upsert(leetcodeUsername, githubUsername string, lastUpdatedAt time.Time) error
}

type LastUpdatedTimestampRepository struct {
	db *pgxpool.Pool
}

func (l LastUpdatedTimestampRepository) Upsert(leetcodeUsername string, githubUsername string, lastUpdatedAt time.Time) error {

	query := `
		INSERT INTO last_updated_timestamp (leetcode_username,github_username , last_updated_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (leetcode_username) DO UPDATE SET
			last_updated_at = $3
	`

	// Execute the upsert command with the parameters
	_, err := l.db.Exec(context.Background(), query, leetcodeUsername, githubUsername, lastUpdatedAt)

	if err != nil {
		logger.Log.Error("LastUpdatedTimestampRepository.Upsert", zap.Error(err))
		return err
	}

	return nil

}

func (l LastUpdatedTimestampRepository) Get(leetcodeUsername string) (*entity.LastUpdatedTimestamp, error) {
	var lastUpdatedTimestamp entity.LastUpdatedTimestamp
	err := l.db.QueryRow(context.Background(), "select * from last_updated_timestamp where leetcode_username = $1", leetcodeUsername).Scan(&lastUpdatedTimestamp.LeetcodeUsername, &lastUpdatedTimestamp.GithubUsername, &lastUpdatedTimestamp.LastUpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		logger.Log.Error("LastUpdatedTimestampRepository.Get", zap.Error(err))
		return nil, err
	}

	return &lastUpdatedTimestamp, nil
}

func NewLastUpdatedTimestampRepository(db *pgxpool.Pool) LastUpdatedTimestampRepositoryInterface {
	return &LastUpdatedTimestampRepository{
		db: db,
	}
}
