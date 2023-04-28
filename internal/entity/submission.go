package entity

import "time"

type CrawlSubmission struct {
	Id               string
	Sha              *string
	UpdatedAt        time.Time
	ProblemId        string
	SubmissionId     string
	LeetcodeUsername string
	GithubUsername   string
}
