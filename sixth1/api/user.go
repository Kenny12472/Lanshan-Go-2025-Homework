package api

import (
	"net/http"
	"time"

	"sixth/dao"
	"sixth/model"
	"sixth/utils"

	"github.com/gin-gonic/gin"
)

// Register 用户注册
func Register(c *gin.Context) {
	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	// 检查用户是否已存在
	if dao.FindUser(req.Username, req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user already exists"})
		return
	}

	// 添加用户到数据库（模拟）
	dao.AddUser(req.Username, req.Password)

	c.JSON(http.StatusOK, gin.H{"message": "register success"})
}

// Login 用户登录
func Login(c *gin.Context) {
	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	// 检查用户是否存在
	if !dao.FindUser(req.Username, req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user not found"})
		return
	}

	// 生成 access token（10分钟）
	token, expire, err := utils.MakeToken(req.Username, time.Now().Add(10*time.Minute))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	// 生成 refresh token（24小时）
	refreshToken, _, _ := utils.MakeRefreshToken(req.Username)

	c.JSON(http.StatusOK, gin.H{
		"message":       "login success",
		"access_token":  token,
		"expire":        expire,
		"refresh_token": refreshToken,
	})
}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	// 从请求头获取 token
	tokenStr := c.GetHeader("Authorization")
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "no token"})
		return
	}

	// 解析 token
	username, _, err := utils.ParseToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
		return
	}

	var req struct {
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	// 修改密码
	dao.UpdatePassword(username, req.NewPassword)

	c.JSON(http.StatusOK, gin.H{"message": "password changed"})
}
