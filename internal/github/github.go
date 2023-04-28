package github

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
	"leetcode-submission-crawler/config"
	"leetcode-submission-crawler/logger"
)

type UseCase interface {
	CreateNewCommit(author, repo, dir, filename, msg, content string, sha *string) error
	GetSHA(username, repo, path string) (*string, error)
}

type Github struct {
	config *config.Config
}

func (g *Github) GetSHA(user, repo, path string) (*string, error) {
	url := fmt.Sprintf(`https://api.github.com/repos/%s/%s/contents/%s`, user, repo, path)

	logger.Log.Info("GetSHA", zap.String("url", url))

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", g.config.Github.Token))

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var githubResponse struct {
		SHA string `json:"sha"`
	}

	err = json.NewDecoder(resp.Body).Decode(&githubResponse)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return &githubResponse.SHA, nil
}

func (g *Github) CreateNewCommit(author, repo, dir, filename, msg, content string, sha *string) error {
	url := fmt.Sprintf(`https://api.github.com/repos/%s/%s/contents/%s/%s`, author, repo, dir, filename)

	logger.Log.Info("CreateNewCommit", zap.String("url", url))

	client := &http.Client{}

	createNewCommitRequest := struct {
		Message string  `json:"message"`
		Content string  `json:"content"`
		Sha     *string `json:"sha,omitempty"`
	}{
		Message: msg,
		Content: base64.StdEncoding.EncodeToString([]byte(content)),
		Sha:     sha,
	}

	reqBody, _ := json.Marshal(createNewCommitRequest)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBody))

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", g.config.Github.Token))

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// read response body
	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d, %s", resp.StatusCode, respBody)
	}

	return nil
}

func NewGithub(config *config.Config) UseCase {
	return &Github{
		config: config,
	}
}
