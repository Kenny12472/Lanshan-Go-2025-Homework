package api

import "github.com/gin-gonic/gin"

// InitRouter_gin 初始化路由
func InitRouter_gin() {
	r := gin.Default()

	// 示例路由
	r.POST("/register", Register)
	r.POST("/login", Login)
	r.POST("/change_password", ChangePassword)

	r.Run(":8080")
}
