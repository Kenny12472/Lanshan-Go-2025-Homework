package main

import (
	"seventh/api"
	"seventh/database"

	"github.com/gin-gonic/gin"
)

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/todo?charset=utf8mb4&parseTime=True&loc=Local"
	if err := database.InitDB(dsn); err != nil {
		panic(err)
	}

	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// 路由
	todos := r.Group("/todos")
	{
		todos.GET("", api.GetTodos)          // 查询全部
		todos.POST("", api.CreateTodo)       // 新增
		todos.PUT("/:id", api.UpdateTodo)    // 修改
		todos.DELETE("/:id", api.DeleteTodo) // 删除
	}

	r.Run(":8080")
}
