package handler

import (
	"net/http"
	"strconv"
	"time"

	"project/internal/db"
	"project/internal/model"

	"github.com/gin-gonic/gin"
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

// helper: 从 gin.Context 安全地读取 user_id（兼容多种类型）
func getUserIDFromContext(c *gin.Context) (uint64, bool) {
	val, ok := c.Get("user_id")
	if !ok || val == nil {
		return 0, false
	}
	switch v := val.(type) {
	case uint64:
		return v, true
	case uint:
		return uint64(v), true
	case int:
		if v < 0 {
			return 0, false
		}
		return uint64(v), true
	case int64:
		if v < 0 {
			return 0, false
		}
		return uint64(v), true
	case float64:
		if v < 0 {
			return 0, false
		}
		return uint64(v), true
	case string:
		// 如果中间件把 id 存为 string
		if id, err := strconv.ParseUint(v, 10, 64); err == nil {
			return id, true
		}
		return 0, false
	default:
		return 0, false
	}
}

// ---------------- ListArticles (分页 + 搜索) ----------------
// 这个版本使用 users.username 作为作者名（兼容之前 DB），并按发布时间倒序返回
func ListArticles(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("page_size", "10")
	query := c.Query("query")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	size, _ := strconv.Atoi(sizeStr)
	if size <= 0 || size > 200 {
		size = 10
	}
	offset := (page - 1) * size

	var results []ArticleResp
	dbQuery := db.DB.Table("articles").
		Select("articles.id, articles.author_id, users.username as author_name, articles.title, articles.content, articles.status, articles.like_count, articles.created_at").
		Joins("left join users on users.id = articles.author_id").
		Where("articles.status = ?", model.ArticlePublish)

	if query != "" {
		likeq := "%" + query + "%"
		dbQuery = dbQuery.Where("(articles.title LIKE ? OR articles.content LIKE ?)", likeq, likeq)
	}

	var total int64
	if err := dbQuery.Order("articles.created_at desc").Count(&total).Offset(offset).Limit(size).Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败", "detail": err.Error()})
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

	var article ArticleResp
	if err := db.DB.Table("articles").
		Select("articles.id, articles.author_id, users.username as author_name, articles.title, articles.content, articles.status, articles.like_count, articles.created_at").
		Joins("left join users on users.id = articles.author_id").
		Where("articles.id = ? AND articles.status = ?", idStr, model.ArticlePublish).
		First(&article).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 增加 view_count（非致命）
	_ = db.DB.Model(&model.Article{}).Where("id = ?", idStr).
		UpdateColumn("view_count", db.DB.Statement.Clauses).Error

	c.JSON(http.StatusOK, article)
}

// ---------------- CreateArticle ----------------
// POST /articles
// body: { title, content, status }
func CreateArticle(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Status  int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误", "detail": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "创建成功", "id": art.ID})
}

// ---------------- UpdateArticle ---------------- (保持不变，简单权限校验)
func UpdateArticle(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	idStr := c.Param("id")
	var article model.Article
	if err := db.DB.First(&article, idStr).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// ---------------- DeleteArticle（逻辑删除） ----------------
func DeleteArticle(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	idStr := c.Param("id")
	var article model.Article
	if err := db.DB.First(&article, idStr).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}
	if article.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限"})
		return
	}

	if err := db.DB.Model(&model.Article{}).Where("id = ?", article.ID).
		Update("status", model.ArticleDeleted).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
