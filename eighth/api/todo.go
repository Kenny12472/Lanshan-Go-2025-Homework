package api

import (
	"encoding/json"
	"net/http"
	"time"

	"eighth/cache"
	"eighth/database"
	"eighth/model"

	"github.com/gin-gonic/gin"
)

func GetTodos(c *gin.Context) {
	cacheKey := "todos:all"
	lockKey := "lock:todos:all"

	val, err := database.RDB.Get(database.Ctx, cacheKey).Result()
	if err == nil {
		var list []model.Todo
		_ = json.Unmarshal([]byte(val), &list)
		c.JSON(http.StatusOK, list)
		return
	}

	if !cache.TryLock(lockKey, time.Second*3) {
		time.Sleep(100 * time.Millisecond)
		c.JSON(http.StatusOK, gin.H{"msg": "请重试"})
		return
	}
	defer cache.Unlock(lockKey)

	var list []model.Todo
	if err := database.DB.Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data, _ := json.Marshal(list)
	database.RDB.Set(database.Ctx, cacheKey, data, cache.RandomExpire())

	c.JSON(http.StatusOK, list)
}

func CreateTodo(c *gin.Context) {
	var todo model.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cache.BloomAdd(todo.ID)
	database.RDB.Del(database.Ctx, "todos:all")

	c.JSON(http.StatusOK, todo)
}

func DeleteTodo(c *gin.Context) {
	id := c.Param("id")

	database.DB.Delete(&model.Todo{}, id)
	database.RDB.Del(database.Ctx, "todos:all")

	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}
