package handler

import (
	"net/http"

	"project/internal/db"
	"project/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ToggleArticleLike(c *gin.Context) {
	uidVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := uidVal.(uint64)
	articleID := c.Param("id")

	var article model.Article
	if err := db.DB.First(&article, articleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 事务：查找是否已点赞，存在则取消，否则新增
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var like model.Like
		res := tx.Where("target_type = ? AND target_id = ? AND user_id = ?", "article", article.ID, userID).First(&like)
		if res.Error == nil {
			// 已点赞，删除
			if err := tx.Delete(&like).Error; err != nil {
				return err
			}
			if err := tx.Model(&model.Article{}).Where("id = ?", article.ID).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error; err != nil {
				return err
			}
			c.JSON(http.StatusOK, gin.H{"liked": false})
			return nil
		}
		if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
			return res.Error
		}
		// 新增点赞
		if err := tx.Create(&model.Like{TargetType: "article", TargetID: article.ID, UserID: userID}).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Article{}).Where("id = ?", article.ID).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error; err != nil {
			return err
		}
		c.JSON(http.StatusOK, gin.H{"liked": true})
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
	}
}

func ToggleCommentLike(c *gin.Context) {
	uidVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := uidVal.(uint64)
	commentID := c.Param("id")

	var comment model.Comment
	if err := db.DB.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var like model.Like
		res := tx.Where("target_type = ? AND target_id = ? AND user_id = ?", "comment", comment.ID, userID).First(&like)
		if res.Error == nil {
			if err := tx.Delete(&like).Error; err != nil {
				return err
			}
			if err := tx.Model(&model.Comment{}).Where("id = ?", comment.ID).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error; err != nil {
				return err
			}
			c.JSON(http.StatusOK, gin.H{"liked": false})
			return nil
		}
		if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
			return res.Error
		}
		if err := tx.Create(&model.Like{TargetType: "comment", TargetID: comment.ID, UserID: userID}).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Comment{}).Where("id = ?", comment.ID).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error; err != nil {
			return err
		}
		c.JSON(http.StatusOK, gin.H{"liked": true})
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
	}
}
