package leetcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"leetcode-submission-crawler/config"
	"leetcode-submission-crawler/logger"
)

func TestMain(m *testing.M) {
	logger.InitLogger()
	NewLeetcodeInstance()
	m.Run()
}

var leetcode UseCase

func NewLeetcodeInstance() {
	var mockConfig = &config.Config{
		Leetcode: config.LeetcodeConfig{
			Username:     "shivaluma",
			CSRFToken:    "pbYak2UyxrRAS9fjrevqiQ2NocQATYfwHsE4rxphXUxaxorPMUUd67HjNax34AbC",
			SessionToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJfYXV0aF91c2VyX2lkIjoiMjY2OTQwNiIsIl9hdXRoX3VzZXJfYmFja2VuZCI6ImFsbGF1dGguYWNjb3VudC5hdXRoX2JhY2tlbmRzLkF1dGhlbnRpY2F0aW9uQmFja2VuZCIsIl9hdXRoX3VzZXJfaGFzaCI6ImRkZGQxODY5Mzk1NzRhZDJiZDhjNjVlMmEzMjI4YTc2MTlhMDU5MjEiLCJpZCI6MjY2OTQwNiwiZW1haWwiOiJzaGl2YWx1bWFAZ21haWwuY29tIiwidXNlcm5hbWUiOiJzaGl2YWx1bWEiLCJ1c2VyX3NsdWciOiJzaGl2YWx1bWEiLCJhdmF0YXIiOiJodHRwczovL2Fzc2V0cy5sZWV0Y29kZS5jb20vdXNlcnMvYXZhdGFycy9hdmF0YXJfMTY3NTUyODI3Ny5wbmciLCJyZWZyZXNoZWRfYXQiOjE2ODI1MzA2MTcsImlwIjoiMTI1LjIxMi4yMjAuMjciLCJpZGVudGl0eSI6Ijg5NGRjNjBhNGUxNDhmNDY1MjYxNWVkMjQ2ZDNlMjk4Iiwic2Vzc2lvbl9pZCI6MzQwNDM1ODd9.H8Rs3cAHthbxpl7yg0cGy97htlAP-lVMCgBM9kyisq8",
		},
	}

	leetcode = NewLeetcode(mockConfig)
}

func TestGetRecentSubmissions(t *testing.T) {
	assert.NotNil(t, leetcode.GetRecentSubmissions())
}

func TestGetSubmissionDetail(t *testing.T) {

	submissionDetail := leetcode.GetSubmissionDetail(940444509)

	assert.NotNil(t, submissionDetail)

	logger.Log.Info("submission detail", zap.Any("submissionDetail", submissionDetail))
}
