package config

import (
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"leetcode-submission-crawler/logger"
)

type Config struct {
	CronExpression string         `mapstructure:"cron_expression"`
	Github         GithubConfig   `mapstructure:"github"`
	Leetcode       LeetcodeConfig `mapstructure:"leetcode"`
	Postgres       PostgresConfig `mapstructure:"postgres"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Driver   string `mapstructure:"driver"`
	Database string `mapstructure:"database"`
}

type LeetcodeConfig struct {
	Username     string `mapstructure:"username"`
	CSRFToken    string `mapstructure:"csrf_token"`
	SessionToken string `mapstructure:"session_token"`
}

type GithubConfig struct {
	Username string `mapstructure:"username"`
	Repo     string `mapstructure:"repo"`
	Token    string `mapstructure:"token"`
}

func GetConfigName() string {
	if len(os.Getenv("CONFIG_NAME")) > 0 {
		return os.Getenv("CONFIG_NAME")
	}
	return "config.local"
}

func GetConfigDirectory() string {
	filePath := os.Getenv("CONFIG_DIRECTORY")
	if filePath != "" {
		return filePath
	}
	return RootDir()
}
func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}

// ReadConfig read config using viper
func ReadConfig() *Config {
	v, err := LoadConfig(GetConfigName(), GetConfigDirectory())
	if err != nil {
		logger.Log.Panic("unable to load config", zap.Error(err))
	}

	config, err := ParseConfig(v)

	if err != nil {
		logger.Log.Panic("unable to parse config", zap.Error(err))
	}
	return config
}

// LoadConfig config file from given path
func LoadConfig(filename, path string) (*viper.Viper, error) {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName(filename)
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return v, nil
}

// ParseConfig file from the given viper
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config
	err := v.Unmarshal(&c)
	if err != nil {
		logger.Log.Fatal("unable to decode into struct", zap.Error(err))
		return nil, err
	}
	return &c, nil
}
