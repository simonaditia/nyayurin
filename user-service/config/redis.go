package config

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		// Addr: "localhost:6379",
		Addr:     "redis-16465.c334.asia-southeast2-1.gce.redns.redis-cloud.com:16465",
		Password: "KHWZnBzAzbfy7Swwvnur1uNYEGGXuagi", // no password set
		DB:       0,                                  // "database-MFP91Y1L"                // use default DB
	})

	_, err := client.Ping(Ctx).Result()
	if err != nil {
		panic(err)
	}
	return client
}
