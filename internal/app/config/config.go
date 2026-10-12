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
	JWT         JWTConfig
	Redis       RedisConfig
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

	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		return nil, fmt.Errorf("critical error: JWT_SECRET_KEY environment variable is not set")
	}

	redisPort, _ := strconv.Atoi(getEnv("REDIS_PORT", "6379"))
	servicePort, _ := strconv.Atoi(getEnv("SERVICE_PORT", "8080"))

	return &Config{
		ServiceHost: getEnv("SERVICE_HOST", "0.0.0.0"),
		ServicePort: servicePort,
		JWT: JWTConfig{
			Token:     jwtSecret,
			ExpiresIn: 24 * time.Hour,
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "0.0.0.0"),
			Port:     redisPort,
			Password: getEnv("REDIS_PASSWORD", "password"),
			User:     getEnv("REDIS_USER", ""),
		},
	}, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
