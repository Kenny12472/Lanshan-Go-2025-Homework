package database

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var (
	RDB *redis.Client
	Ctx = context.Background()
)

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	if err := RDB.Ping(Ctx).Err(); err != nil {
		log.Fatalf("Redis连接失败: %v", err)
	}
}
