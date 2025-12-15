package cache

import (
	"time"

	"eighth/database"
)

func TryLock(key string, ttl time.Duration) bool {
	ok, _ := database.RDB.SetNX(database.Ctx, key, 1, ttl).Result()
	return ok
}

func Unlock(key string) {
	database.RDB.Del(database.Ctx, key)
}
