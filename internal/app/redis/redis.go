package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

const servicePrefix = "composers_service."

type Client struct {
	client *redis.Client
}

func New(ctx context.Context, host string, port int, password string, user string) (*Client, error) {
	rClient := redis.NewClient(&redis.Options{
		Addr:        host + ":" + strconv.Itoa(port),
		Username:    user,
		Password:    password,
		DB:          0,
		DialTimeout: 10 * time.Second,
		ReadTimeout: 10 * time.Second,
	})

	if _, err := rClient.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("cant ping redis: %w", err)
	}

	return &Client{client: rClient}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

const jwtPrefix = "jwt."

func getJWTKey(token string) string {
	return servicePrefix + jwtPrefix + token
}

func (c *Client) WriteJWTToBlacklist(ctx context.Context, jwtStr string, jwtTTL time.Duration) error {
	return c.client.Set(ctx, getJWTKey(jwtStr), true, jwtTTL).Err()
}

func (c *Client) CheckJWTInBlacklist(ctx context.Context, jwtStr string) error {
	return c.client.Get(ctx, getJWTKey(jwtStr)).Err()
}
