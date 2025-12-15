package main

import (
	"eighth/api"
	"eighth/database"

	"github.com/gin-gonic/gin"
)

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/eighth?charset=utf8mb4&parseTime=True&loc=Local"
	if err := database.InitDB(dsn); err != nil {
		panic(err)
	}

	database.InitRedis()

	r := gin.Default()

	r.SetTrustedProxies([]string{"127.0.0.1"})

	todos := r.Group("/todos")
	{
		todos.GET("", api.GetTodos)
		todos.POST("", api.CreateTodo)
	}

	r.Run(":8080")
}
