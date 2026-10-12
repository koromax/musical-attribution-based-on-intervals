package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServiceHost string
	ServicePort int
	Database    DatabaseConfig
	Minio       MinioConfig
	JWT         JWTConfig
	Redis       RedisConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
}

type JWTConfig struct {
	Token     string
	ExpiresIn time.Duration
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	User     string
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load()

	redisPort, err := getPort("REDIS_PORT", "6379")
	if err != nil {
		return nil, err
	}
	servicePort, err := getPort("SERVICE_PORT", "8080")
	if err != nil {
		return nil, err
	}
	databasePort, err := getPort("DB_PORT", "5432")
	if err != nil {
		return nil, err
	}

	return &Config{
		ServiceHost: getEnv("SERVICE_HOST", "0.0.0.0"),
		ServicePort: servicePort,
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     databasePort,
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASS"),
			Name:     os.Getenv("DB_NAME"),
		},
		Minio: MinioConfig{
			Endpoint:  os.Getenv("MINIO_ENDPOINT"),
			AccessKey: os.Getenv("MINIO_ACCESS_KEY"),
			SecretKey: os.Getenv("MINIO_SECRET_KEY"),
			Bucket:    os.Getenv("MINIO_BUCKET_NAME"),
		},
		JWT: JWTConfig{
			Token:     os.Getenv("JWT_SECRET_KEY"),
			ExpiresIn: 24 * time.Hour,
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     redisPort,
			Password: getEnv("REDIS_PASSWORD", "password"),
			User:     getEnv("REDIS_USER", ""),
		},
	}, nil
}

func getPort(key, defaultVal string) (int, error) {
	port, err := strconv.Atoi(getEnv(key, defaultVal))
	if err != nil || port < 1 || port > 65535 {
		return 0, fmt.Errorf("%s must be a valid port number", key)
	}
	return port, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
