package config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var config *Config

type Config struct {
	Database DBConfig     `mapstructure:"database"`
	Server   ServerConfig `mapstructure:"server"`
	Redis    RedisConfig  `mapstructure:"redis"`
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

func setDefaults() {
	viper.SetDefault("database.name", "trafficguard_db")
	viper.SetDefault("database.username", "trafficguard")
	viper.SetDefault("database.password", "password") // kendi şifreniz
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5440")

	viper.SetDefault("server.port", "8080")

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
