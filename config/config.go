package config

type Config struct {
	CronExpression string         `mapstructure:"cron_expression"`
	Github         GithubConfig   `mapstructure:"github"`
	Leetcode       LeetcodeConfig `mapstructure:"leetcode"`
}

type LeetcodeConfig struct {
	Username     string `mapstructure:"username"`
	CSRFToken    string `mapstructure:"csrf_token"`
	SessionToken string `mapstructure:"session_token"`
}

type GithubConfig struct {
	Username string `mapstructure:"username"`
}
