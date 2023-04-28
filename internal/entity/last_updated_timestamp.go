package entity

import "time"

type LastUpdatedTimestamp struct {
	LeetcodeUsername string    `json:"leetcode_username"`
	GithubUsername   string    `json:"github_username"`
	LastUpdatedAt    time.Time `json:"last_updated_at"`
}
