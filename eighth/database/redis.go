package database

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var (
	RDB *redis.Client
	Ctx = context.Background()
)

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})

	if _, err := RDB.Ping(Ctx).Result(); err != nil {
		panic(err)
	}
}
