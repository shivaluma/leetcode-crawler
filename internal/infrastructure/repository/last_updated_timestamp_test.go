package repository

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"leetcode-submission-crawler/config"
	"leetcode-submission-crawler/database"
)

var db *pgxpool.Pool
var repo LastUpdatedTimestampRepositoryInterface

func TestMain(m *testing.M) {

	config := config.ReadConfig()

	db = database.NewDBClient(config)

	m.Run()
}

func TestNewLastUpdatedTimestampRepository(t *testing.T) {
	repo = NewLastUpdatedTimestampRepository(db)
}

func TestLastUpdatedTimestampRepository_Upsert(t *testing.T) {
	err := repo.Upsert("leetcode", "github", time.Now())
	assert.Nil(t, err)

	data, err := repo.Get("leetcode")
	assert.Nil(t, err)

	assert.Equal(t, "leetcode", data.LeetcodeUsername)
}

func TestLastUpdatedTimestampRepository(t *testing.T) {
	data, err := repo.Get("leetcode")
	assert.Nil(t, err)

	assert.Equal(t, "leetcode", data.LeetcodeUsername)

	ts := data.LastUpdatedAt

	err = repo.Upsert("leetcode", "github", time.Now())
	assert.Nil(t, err)

	newData, err := repo.Get("leetcode")

	assert.Equal(t, data.GithubUsername, newData.GithubUsername)
	assert.True(t, ts.Before(newData.LastUpdatedAt))

}
