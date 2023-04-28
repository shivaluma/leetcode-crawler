package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"leetcode-submission-crawler/config"
	"leetcode-submission-crawler/logger"
)

func NewDBClient(config *config.Config) *pgxpool.Pool {
	db, err := pgxpool.New(context.Background(), PostgresConnectionString(config.Postgres))

	if err != nil {
		logger.Log.Panic("unable to connect to database", zap.Error(err))
	}

	return db
}

func PostgresConnectionString(config config.PostgresConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.Username, config.Password, config.Host, config.Port, config.Database)
}
