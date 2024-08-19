package redis

import (
	"github.com/go-redis/redis/v8"
	_ "github.com/go-redis/redis/v8"
)

func RedisConnection() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "root",
		DB:       0,
	})

	return rdb

}
