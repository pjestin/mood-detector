package io

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client redis.Client
	ctx    context.Context
}

const REDIS_ADDRESS = "redis:6379"

func (c *RedisClient) Init(db int) {
	c.client = *redis.NewClient(&redis.Options{
		Addr:     REDIS_ADDRESS,
		Password: "", // no password set
		DB:       db,
	})
	c.ctx = context.Background()
}

func (c *RedisClient) Set(key string, value interface{}) error {
	return c.client.Set(c.ctx, key, value, 0).Err()
}

func (c *RedisClient) Get(key string) (string, error) {
	res := c.client.Get(c.ctx, key)
	if res.Err() != nil {
		return "", res.Err()
	}

	return res.String(), nil
}
