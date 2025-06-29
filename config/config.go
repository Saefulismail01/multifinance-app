package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Driver   string
}

type APIConfig struct {
	ApiPort string
}

type Config struct {
	DBConfig
	APIConfig
}

func (c *Config) readConfig() error {
	_ = godotenv.Load()

	// Set default values
	c.DBConfig = DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3306"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "multifinance"),
		Driver:   getEnv("DB_DRIVER", "mysql"),
	}

	c.APIConfig = APIConfig{
		ApiPort: getEnv("API_PORT", "8080"),
	}


	if c.DBConfig.Host == "" || c.DBConfig.Port == "" || 
	   c.DBConfig.User == "" || c.DBConfig.DBName == "" || 
	   c.APIConfig.ApiPort == "" {
		return fmt.Errorf("missing required configuration")
	}


	if _, err := strconv.Atoi(c.DBConfig.Port); err != nil {
		return fmt.Errorf("invalid DB port number: %v", err)
	}

	if _, err := strconv.Atoi(c.APIConfig.ApiPort); err != nil {
		return fmt.Errorf("invalid API port number: %v", err)
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := cfg.readConfig(); err != nil {
		return nil, err
	}
	return cfg, nil
}

