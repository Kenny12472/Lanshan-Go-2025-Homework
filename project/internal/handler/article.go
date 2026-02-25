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

// ��ѧ��ע�ͣ� ArticleResp ���ڷ��ظ�ǰ�˵����½ṹ����������??
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

// ��ѧ��ע�ͣ� ---------------- ListArticles (��ҳ + ���� + ��ע����) ----------------
// ��ѧ��ע�ͣ� GET /articles?page=1&page_size=10&query=xxx&followed_first=1
func ListArticles(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("page_size", "10")
	query := c.Query("query")
	followedFirst := c.Query("followed_first") // "1" ??"true" ��������

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

	// ��ѧ��ע�ͣ� ���������Ҫ��ע���ȣ������е�¼�û������� left join follows f ���� order ��������ʾ��ע��??
	if followedFirst == "1" || followedFirst == "true" {
		if uidVal, ok := c.Get("user_id"); ok {
			userID := uidVal.(uint64)
			dbQuery = dbQuery.
				Joins("LEFT JOIN follows f ON f.following_id = articles.author_id AND f.follower_id = ?", userID).
				Order("CASE WHEN f.id IS NULL THEN 1 ELSE 0 END, articles.created_at DESC")
		} else {
			// ��ѧ��ע�ͣ� δ��¼�Ͱ�ʱ�䵹��
			dbQuery = dbQuery.Order("articles.created_at desc")
		}
	} else {
		// ��ѧ��ע�ͣ� Ĭ�ϰ�����ʱ�䵹������ǰ??
		dbQuery = dbQuery.Order("articles.created_at desc")
	}

	var total int64
	dbQuery.Count(&total)

	if err := dbQuery.Offset(offset).Limit(size).Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":      page,
		"page_size": size,
		"total":     total,
		"items":     results,
	})
}

// ��ѧ��ע�ͣ� ---------------- GetArticle (��ƪ����������) ----------------
func GetArticle(c *gin.Context) {
	idStr := c.Param("id")

	// ��ѧ��ע�ͣ� ��ѯ���£�������״̬��
	var article ArticleResp
	if err := db.DB.Table("articles").
		Select("articles.id, articles.author_id, COALESCE(users.display_name, users.username) as author_name, articles.title, articles.content, articles.status, articles.like_count, articles.created_at").
		Joins("left join users on users.id = articles.author_id").
		Where("articles.id = ? AND articles.status = ?", idStr, model.ArticlePublish).
		First(&article).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	// ��ѧ��ע�ͣ� ���� view_count��ʹ??gorm.Expr??
	_ = db.DB.Model(&model.Article{}).Where("id = ?", idStr).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error

	c.JSON(http.StatusOK, article)
}

// ��ѧ��ע�ͣ� ---------------- CreateArticle ----------------
// ��ѧ��ע�ͣ� POST /articles
// ��ѧ��ע�ͣ� body: { title, content, status }
func CreateArticle(c *gin.Context) {
	uidVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := uidVal.(uint64)

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Status  int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	if req.Title == "" || req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title/content required"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "created", "id": art.ID})
}

// ��ѧ��ע�ͣ� ---------------- UpdateArticle ----------------
// ��ѧ��ע�ͣ� PUT /articles/:id
// ��ѧ��ע�ͣ� body: { title, content, status }
func UpdateArticle(c *gin.Context) {
	uidVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := uidVal.(uint64)

	idStr := c.Param("id")
	var article model.Article
	if err := db.DB.First(&article, idStr).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	// ��ѧ��ע�ͣ� ֻ�����߿��Ը��£�Ҳ������չΪ����ԱȨ�ޣ�
	if article.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Status  int    `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	article.Title = req.Title
	article.Content = req.Content
	article.Status = req.Status
	article.UpdatedAt = time.Now()

	if err := db.DB.Save(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// ��ѧ��ע�ͣ� ---------------- DeleteArticle���߼�ɾ��??----------------
// ��ѧ��ע�ͣ� DELETE /articles/:id
func DeleteArticle(c *gin.Context) {
	uidVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := uidVal.(uint64)

	idStr := c.Param("id")
	var article model.Article
	if err := db.DB.First(&article, idStr).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	// ��ѧ��ע�ͣ� ֻ�����߿���ɾ���������Ա??
	if article.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if err := db.DB.Model(&model.Article{}).Where("id = ?", article.ID).
		Update("status", model.ArticleDeleted).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
