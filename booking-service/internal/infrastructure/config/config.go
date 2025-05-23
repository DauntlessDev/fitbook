package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Logging  LoggingConfig
	App      AppConfig
}

type ServerConfig struct {
	Port         int
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	URL string
}

type RedisConfig struct {
	URL string
}

type LoggingConfig struct {
	Level  string
	Format string
}

type AppConfig struct {
	Env         string
	ServiceName string
}

func Load() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Warning: .env file not found or could not be loaded")
	}

	portStr, err := mustGetEnv("BOOKING_SERVER_PORT")
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid port number: %w", err)
	}

	readTimeoutStr, err := mustGetEnv("BOOKING_SERVER_READ_TIMEOUT")
	if err != nil {
		return nil, err
	}
	readTimeout, err := time.ParseDuration(readTimeoutStr)
	if err != nil {
		return nil, fmt.Errorf("invalid read timeout: %w", err)
	}

	writeTimeoutStr, err := mustGetEnv("BOOKING_SERVER_WRITE_TIMEOUT")
	if err != nil {
		return nil, err
	}
	writeTimeout, err := time.ParseDuration(writeTimeoutStr)
	if err != nil {
		return nil, fmt.Errorf("invalid write timeout: %w", err)
	}

	host, err := mustGetEnv("BOOKING_SERVER_HOST")
	if err != nil {
		return nil, err
	}

	dbURL, err := mustGetEnv("DB_URL")
	if err != nil {
		return nil, err
	}

	redisURL, err := mustGetEnv("REDIS_URL")
	if err != nil {
		return nil, err
	}

	logLevel, err := mustGetEnv("BOOKING_LOG_LEVEL")
	if err != nil {
		return nil, err
	}

	logFormat, err := mustGetEnv("BOOKING_LOG_FORMAT")
	if err != nil {
		return nil, err
	}

	env, err := mustGetEnv("BOOKING_ENV")
	if err != nil {
		return nil, err
	}

	serviceName, err := mustGetEnv("BOOKING_SERVICE_NAME")
	if err != nil {
		return nil, err
	}

	return &Config{
		Server: ServerConfig{
			Port:         port,
			Host:         host,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
		Database: DatabaseConfig{
			URL: dbURL,
		},
		Redis: RedisConfig{
			URL: redisURL,
		},
		Logging: LoggingConfig{
			Level:  logLevel,
			Format: logFormat,
		},
		App: AppConfig{
			Env:         env,
			ServiceName: serviceName,
		},
	}, nil
}

func mustGetEnv(key string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}
	return "", fmt.Errorf("environment variable %s is required but not set", key)
}

func (config *DatabaseConfig) GetDSN() string {
	return config.URL
}

func (config *RedisConfig) GetRedisAddr() string {
	return config.URL
}
