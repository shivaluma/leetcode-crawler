package leetcode

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hasura/go-graphql-client"
	"go.uber.org/zap"
	"leetcode-submission-crawler/config"
	"leetcode-submission-crawler/internal/entity"
	"leetcode-submission-crawler/logger"
)

type UseCase interface {
	GetRecentSubmissions(size int) []entity.RecentAcSubmission
	GetSubmissionDetail(id int64) *entity.Submission
	GetQuestionContent(slug string) *entity.QuestionContent
}

func NewLeetcode(config *config.Config) UseCase {
	httpClient := &http.Client{}

	cookies := cookieHelper{}

	cookies.AddCookie(
		&http.Cookie{
			Name:  "csrftoken",
			Value: config.Leetcode.CSRFToken,
		},
		&http.Cookie{
			Name:  "LEETCODE_SESSION",
			Value: config.Leetcode.SessionToken,
		},
	)

	httpClient.Transport = &headerTransport{
		http.DefaultTransport,
		map[string]string{
			"x-csrftoken": config.Leetcode.CSRFToken,
			"Cookie":      cookies.String(),
		},
	}

	return &Impl{
		gc:     graphql.NewClient("https://leetcode.com/graphql", httpClient),
		config: config,
	}
}

type Impl struct {
	gc     *graphql.Client
	config *config.Config
}

func (l *Impl) GetQuestionContent(slug string) *entity.QuestionContent {
	var q struct {
		Question entity.QuestionContent `graphql:"question(titleSlug: $titleSlug)"`
	}

	variables := map[string]interface{}{
		"titleSlug": slug,
	}

	err := l.gc.Query(context.Background(), &q, variables)

	if err != nil {
		logger.Log.Error("leetcode.GetQuestionContent", zap.Error(err))
		return nil
	}

	return &q.Question
}

func (l *Impl) GetRecentSubmissions(limit int) []entity.RecentAcSubmission {

	var q struct {
		RecentAcSubmissionList []entity.RecentAcSubmission `graphql:"recentAcSubmissionList(username: $username, limit: $limit)"`
	}

	variables := map[string]interface{}{
		"username": l.config.Leetcode.Username,
		"limit":    limit,
	}

	err := l.gc.Query(context.Background(), &q, variables)
	if err != nil {
		logger.Log.Error("leetcode.GetRecentSubmissions", zap.Error(err))
		return nil
	}

	return q.RecentAcSubmissionList
}

func (l *Impl) GetSubmissionDetail(id int64) *entity.Submission {

	var q struct {
		SubmissionDetails entity.Submission `graphql:"submissionDetails(submissionId: $submissionId)"`
	}

	variables := map[string]interface{}{
		"submissionId": id,
	}

	err := l.gc.Query(context.Background(), &q, variables)
	if err != nil {
		logger.Log.Error("leetcode.GetSubmissionDetail", zap.Error(err))
		return nil
	}

	return &q.SubmissionDetails
}

// Custom transport that adds headers to each request
type headerTransport struct {
	transport http.RoundTripper
	headers   map[string]string
}

func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, value := range t.headers {
		req.Header.Add(key, value)
	}
	return t.transport.RoundTrip(req)
}

type cookieHelper struct {
	cookies []*http.Cookie
}

func (c *cookieHelper) AddCookie(cookie ...*http.Cookie) {
	c.cookies = append(c.cookies, cookie...)
}

var cookieNameSanitizer = strings.NewReplacer("\n", "-", "\r", "-")

func sanitizeCookieName(n string) string {
	return cookieNameSanitizer.Replace(n)
}

func (c *cookieHelper) String() string {
	sb := &strings.Builder{}
	for _, cookie := range c.cookies {
		sb.WriteString(fmt.Sprintf("%s=%s; ", sanitizeCookieName(cookie.Name), cookie.Value))
	}
	return sb.String()
}

func ExtensionNameMapping(name string) string {
	switch name {
	case "golang":
		return "go"
	case "java":
		return "java"
	case "javascript":
		return "js"
	case "ruby":
		return "rb"
	default:
		return name
	}
}
