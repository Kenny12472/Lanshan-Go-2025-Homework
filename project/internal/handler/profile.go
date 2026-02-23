package handler

import (
	"net/http"
	"strconv"
	"time"

	"project/internal/db"
	"project/internal/model"

	"github.com/gin-gonic/gin"
)

// ProfileResp: 返回给前端的用户信息
type ProfileResp struct {
	ID          uint64    `json:"id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name"`
	Bio         string    `json:"bio"`
	CreatedAt   time.Time `json:"created_at"`
}

// GetMyProfile - GET /me/profile
func GetMyProfile(c *gin.Context) {
	uidVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := uidVal.(uint64)

	var u model.User
	if err := db.DB.First(&u, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	resp := ProfileResp{
		ID:          u.ID,
		Username:    u.Username,
		DisplayName: u.DisplayName,
		Bio:         u.Bio,
		CreatedAt:   u.CreatedAt,
	}
	c.JSON(http.StatusOK, resp)
}

// UpdateMyProfile - PUT /me/profile
// body: { display_name, bio }
func UpdateMyProfile(c *gin.Context) {
	uidVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := uidVal.(uint64)

	var req struct {
		DisplayName string `json:"display_name"`
		Bio         string `json:"bio"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	// 简单校验 (可加强)
	if len(req.DisplayName) > 128 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "显示名太长"})
		return
	}

	if err := db.DB.Model(&model.User{}).Where("id = ?", userID).
		Updates(map[string]interface{}{"display_name": req.DisplayName, "bio": req.Bio, "updated_at": time.Now()}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// GetUserProfile - GET /users/:id/profile （公开）
func GetUserProfile(c *gin.Context) {
	idStr := c.Param("id")
	uid, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID错误"})
		return
	}
	var u model.User
	if err := db.DB.First(&u, uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	resp := ProfileResp{ID: u.ID, Username: u.Username, DisplayName: u.DisplayName, Bio: u.Bio, CreatedAt: u.CreatedAt}
	c.JSON(http.StatusOK, resp)
}
