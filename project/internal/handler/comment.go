package handler

import (
	"net/http"
	"time"

	"project/internal/db"
	"project/internal/model"

	"github.com/gin-gonic/gin"
)

type CommentResp struct {
	ID        uint64    `json:"id"`
	ArticleID uint64    `json:"article_id"`
	UserID    uint64    `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	LikeCount int64     `json:"like_count"`
	CreatedAt time.Time `json:"created_at"`
}

func PostComment(c *gin.Context) {
	uidVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := uidVal.(uint64)

	articleIDStr := c.Param("id")
	var req struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 转换 articleID
	var article model.Article
	if err := db.DB.First(&article, articleIDStr).Error; err != nil || article.Status != model.ArticlePublish {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	comment := model.Comment{
		ArticleID: article.ID,
		UserID:    userID,
		Content:   req.Content,
	}
	if err := db.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	// 返回创建的数据（带用户名）
	var resp CommentResp
	db.DB.Table("comments").
		Select("comments.id, comments.article_id, comments.user_id, users.username, comments.content, comments.like_count, comments.created_at").
		Joins("left join users on users.id = comments.user_id").
		Where("comments.id = ?", comment.ID).
		Scan(&resp)

	c.JSON(http.StatusOK, resp)
}

func ListComments(c *gin.Context) {
	articleID := c.Param("id")
	var list []CommentResp
	db.DB.Table("comments").
		Select("comments.id, comments.article_id, comments.user_id, users.username, comments.content, comments.like_count, comments.created_at").
		Joins("left join users on users.id = comments.user_id").
		Where("comments.article_id = ?", articleID).
		Order("comments.created_at desc").
		Find(&list)
	c.JSON(http.StatusOK, list)
}
