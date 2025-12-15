package cache

import "eighth/database"

const BloomKey = "bf:todos"

func BloomAdd(id uint) {
	database.RDB.SAdd(database.Ctx, BloomKey, id)
}

func BloomExists(id uint) bool {
	ok, _ := database.RDB.SIsMember(database.Ctx, BloomKey, id).Result()
	return ok
}
