package handler

import (
	"net/http"
	"strconv"

	"project/internal/db"
	"project/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ToggleFollow(c *gin.Context) {
	uidVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := uidVal.(uint64)

	targetStr := c.Param("id")
	targetID64, err := strconv.ParseUint(targetStr, 10, 64)
	if err != nil || targetID64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID错误"})
		return
	}
	if targetID64 == userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能关注自己"})
		return
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		var f model.Follow
		res := tx.Where("follower_id = ? AND following_id = ?", userID, targetID64).First(&f)
		if res.Error == nil {
			if err := tx.Delete(&f).Error; err != nil {
				return err
			}
			c.JSON(http.StatusOK, gin.H{"following": false})
			return nil
		}
		if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
			return res.Error
		}
		if err := tx.Create(&model.Follow{FollowerID: userID, FollowingID: targetID64}).Error; err != nil {
			return err
		}
		c.JSON(http.StatusOK, gin.H{"following": true})
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
	}
}

func GetMyFollows(c *gin.Context) {
	uidVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := uidVal.(uint64)

	var ids []uint64
	db.DB.Model(&model.Follow{}).Where("follower_id = ?", userID).Pluck("following_id", &ids)
	c.JSON(http.StatusOK, gin.H{"followings": ids})
}
