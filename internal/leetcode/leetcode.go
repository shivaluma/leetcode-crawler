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

type Leetcode interface {
	GetRecentSubmissions() []entity.RecentAcSubmission
	GetSubmissionDetail(id int64) *entity.Submission
}

func NewLeetcode(config *config.Config) Leetcode {
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

func (l *Impl) GetRecentSubmissions() []entity.RecentAcSubmission {

	var q struct {
		RecentAcSubmissionList []entity.RecentAcSubmission `graphql:"recentAcSubmissionList(username: $username, limit: $limit)"`
	}

	variables := map[string]interface{}{
		"username": l.config.Leetcode.Username,
		"limit":    20,
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
		SubmissionDetail entity.Submission `graphql:"submissionDetails(submissionId: $submissionId)"`
	}

	variables := map[string]interface{}{
		"submissionId": id,
	}

	err := l.gc.Query(context.Background(), &q, variables)
	if err != nil {
		logger.Log.Error("leetcode.GetSubmissionDetail", zap.Error(err))
		return nil
	}

	return &q.SubmissionDetail
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
