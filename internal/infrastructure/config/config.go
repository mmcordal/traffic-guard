package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var config *Config

type Config struct {
	Database DBConfig      `mapstructure:"database"`
	Server   ServerConfig  `mapstructure:"server"`
	Redis    RedisConfig   `mapstructure:"redis"`
	Analyze  AnalyzeConfig `mapstructure:"analyze"`
}

type DBConfig struct {
	Name     string `mapstructure:"name"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type AnalyzeConfig struct {
	BucketWindowSeconds int `mapstructure:"bucket_window_seconds"`
	AnalyzeEverySeconds int `mapstructure:"analyze_every_seconds"`
	HistoryBuckets      int `mapstructure:"history_buckets"`
}

func (c AnalyzeConfig) BucketWindow() time.Duration {
	if c.BucketWindowSeconds <= 0 {
		return time.Minute
	}

	return time.Duration(c.BucketWindowSeconds) * time.Second
}

func (c AnalyzeConfig) AnalyzeEvery() time.Duration {
	if c.AnalyzeEverySeconds <= 0 {
		return time.Minute
	}

	return time.Duration(c.AnalyzeEverySeconds) * time.Second
}

func (c AnalyzeConfig) HistoryLimit() int64 {
	if c.HistoryBuckets <= 0 {
		return 10
	}

	return int64(c.HistoryBuckets)
}

func (c AnalyzeConfig) BucketTTL() time.Duration {
	return c.BucketWindow() * time.Duration(c.HistoryLimit()+5)
}

func setDefaults() {
	viper.SetDefault("database.name", "trafficguard_db")
	viper.SetDefault("database.username", "trafficguard")
	viper.SetDefault("database.password", "password") // kendi şifreniz
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5440")

	viper.SetDefault("server.port", "8080")

	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", "6380")
	// prod için
	viper.SetDefault("analyze.bucket_window_seconds", 30) // 60
	viper.SetDefault("analyze.analyze_every_seconds", 30) // 60
	viper.SetDefault("analyze.history_buckets", 10)       // 10
}
func Setup() (*Config, error) {
	setDefaults()

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error loading .env file: %v, loading environment variables instead.", err)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if config == nil {
		config = &Config{}
	}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	if config.Server.Port == "" {
		if p := os.Getenv("SERVER_PORT"); p != "" {
			config.Server.Port = p
		} else {
			config.Server.Port = "8080"
		}
	}

	return config, nil
}

func Get() *Config {
	if config == nil {
		panic("Conifg gelemedi")
	}

	return config
}
