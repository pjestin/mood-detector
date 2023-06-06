package io

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client redis.Client
	ctx    context.Context
}

func (c *RedisClient) Init(db int, redisAddress string) {
	c.client = *redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: "", // no password set
		DB:       db,
	})
	c.ctx = context.Background()
}

func (c *RedisClient) Set(key string, value interface{}) error {
	return c.client.Set(c.ctx, key, value, 0).Err()
}

func (c *RedisClient) GetKeys() ([]string, error) {
	res := c.client.Scan(c.ctx, 0, "*", 0)
	if res.Err() != nil {
		return nil, res.Err()
	}

	iter := res.Iterator()

	var keys []string
	for iter.Next(c.ctx) {
		keys = append(keys, iter.Val())
	}

	return keys, nil
}

func (c *RedisClient) GetEntries() (map[string]string, error) {
	keys, err := c.GetKeys()
	if err != nil {
		return nil, err
	}

	entries := make(map[string]string)

	for _, key := range keys {
		res := c.client.Get(c.ctx, key)
		if res.Err() != nil {
			return nil, res.Err()
		}

		entries[key] = res.Val()
	}

	return entries, nil
}
