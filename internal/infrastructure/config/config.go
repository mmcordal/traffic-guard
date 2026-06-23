package config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var config *Config

type Config struct {
	Database DBConfig
	Server   ServerConfig
	Redis    RedisConfig
}

type DBConfig struct {
	Name     string
	Username string
	Password string
	Host     string
	Port     string
}

type ServerConfig struct {
	Port string
}

type RedisConfig struct {
	Host string
	Port string
}

func setDefaults() {
	viper.SetDefault("database.name", "trafficguard_db")
	viper.SetDefault("database.username", "trafficguard")
	viper.SetDefault("database.password", "trafficguard123") // kendi şifreniz
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5433")

	viper.SetDefault("server.port", "3000")

	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", "6380")
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
			config.Server.Port = "3000"
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
