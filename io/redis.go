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

func (c *RedisClient) Init() {
	c.client = *redis.NewClient(&redis.Options{
		Addr:     REDIS_ADDRESS,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	c.ctx = context.Background()
}

func (c *RedisClient) Set(key string, value int) error {
	return c.client.Set(c.ctx, key, value, 0).Err()
}
