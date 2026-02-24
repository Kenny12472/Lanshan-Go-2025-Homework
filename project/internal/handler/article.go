package handler

import (
	"net/http"
	"strconv"
	"time"

	"project/internal/db"
	"project/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ArticleResp 用于返回给前端的文章结构（含作者名）
type ArticleResp struct {
	ID         uint64    `json:"id"`
	AuthorID   uint64    `json:"author_id"`
	AuthorName string    `json:"author_name"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Status     int       `json:"status"`
	LikeCount  int64     `json:"like_count"`
	CreatedAt  time.Time `json:"created_at"`
}

// ---------------- ListArticles (分页 + 搜索 + 关注优先) ----------------
// GET /articles?page=1&page_size=10&query=xxx&followed_first=1
func ListArticles(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("page_size", "10")
	query := c.Query("query")
	followedFirst := c.Query("followed_first") // "1" 或 "true" 触发优先

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	size, _ := strconv.Atoi(sizeStr)
	if size <= 0 || size > 100 {
		size = 10
	}
	offset := (page - 1) * size

	var results []ArticleResp
	dbQuery := db.DB.Table("articles").
		Select("articles.id, articles.author_id, COALESCE(users.display_name, users.username) as author_name, articles.title, articles.content, articles.status, articles.like_count, articles.created_at").
		Joins("left join users on users.id = articles.author_id").
		Where("articles.status = ?", model.ArticlePublish)

	if query != "" {
		likeq := "%" + query + "%"
		dbQuery = dbQuery.Where("(articles.title LIKE ? OR articles.content LIKE ?)", likeq, likeq)
	}

	// 如果请求需要关注优先，并且有登录用户，则做 left join follows f 并在 order 中优先显示关注作者
	if followedFirst == "1" || followedFirst == "true" {
		if uidVal, ok := c.Get("user_id"); ok {
			userID := uidVal.(uint64)
			dbQuery = dbQuery.
				Joins("LEFT JOIN follows f ON f.following_id = articles.author_id AND f.follower_id = ?", userID).
				Order("CASE WHEN f.id IS NULL THEN 1 ELSE 0 END, articles.created_at DESC")
		} else {
			// 未登录就按时间倒序
			dbQuery = dbQuery.Order("articles.created_at desc")
		}
	} else {
		// 默认按发布时间倒序（最新前）
		dbQuery = dbQuery.Order("articles.created_at desc")
	}

	var total int64
	dbQuery.Count(&total)

	if err := dbQuery.Offset(offset).Limit(size).Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":      page,
		"page_size": size,
		"total":     total,
		"items":     results,
	})
}

// ---------------- GetArticle (单篇，含作者名) ----------------
func GetArticle(c *gin.Context) {
	idStr := c.Param("id")

	// 查询文章（仅发布状态）
	var article ArticleResp
	if err := db.DB.Table("articles").
		Select("articles.id, articles.author_id, COALESCE(users.display_name, users.username) as author_name, articles.title, articles.content, articles.status, articles.like_count, articles.created_at").
		Joins("left join users on users.id = articles.author_id").
		Where("articles.id = ? AND articles.status = ?", idStr, model.ArticlePublish).
		First(&article).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 增加 view_count（使用 gorm.Expr）
	if err := db.DB.Model(&model.Article{}).Where("id = ?", idStr).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error; err != nil {
		// 非致命
	}

	c.JSON(http.StatusOK, article)
}

// ---------------- CreateArticle ----------------
// POST /articles
// body: { title, content, status }
func CreateArticle(c *gin.Context) {
	uidVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := uidVal.(uint64)

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Status  int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if req.Title == "" || req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标题或内容不能为空"})
		return
	}

	art := model.Article{
		AuthorID:  userID,
		Title:     req.Title,
		Content:   req.Content,
		Status:    req.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.DB.Create(&art).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "创建成功", "id": art.ID})
}

// ---------------- UpdateArticle ----------------
// PUT /articles/:id
// body: { title, content, status }
func UpdateArticle(c *gin.Context) {
	uidVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := uidVal.(uint64)

	idStr := c.Param("id")
	var article model.Article
	if err := db.DB.First(&article, idStr).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 只有作者可以更新（也可以扩展为管理员权限）
	if article.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限"})
		return
	}

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Status  int    `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	article.Title = req.Title
	article.Content = req.Content
	article.Status = req.Status
	article.UpdatedAt = time.Now()

	if err := db.DB.Save(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// ---------------- DeleteArticle（逻辑删除） ----------------
// DELETE /articles/:id
func DeleteArticle(c *gin.Context) {
	uidVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID := uidVal.(uint64)

	idStr := c.Param("id")
	var article model.Article
	if err := db.DB.First(&article, idStr).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 只有作者可以删除（或管理员）
	if article.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限"})
		return
	}

	if err := db.DB.Model(&model.Article{}).Where("id = ?", article.ID).
		Update("status", model.ArticleDeleted).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
