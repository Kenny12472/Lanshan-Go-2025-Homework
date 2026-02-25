package handler

import (
	"net/http"
	"time"

	"project/internal/db"
	"project/internal/middleware"
	"project/internal/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ��ѧ��ע�ͣ� Register ʾ����POST /register??
func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "��������"})
		return
	}

	// ��ѧ��ע�ͣ� ����û����Ƿ��Ѵ�??
	var exist model.User
	if err := db.DB.Where("username = ?", req.Username).First(&exist).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "�û����Ѵ���"})
		return
	}

	// ��ѧ��ע�ͣ� ���� bcrypt ��ϣ
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "���봦��ʧ��"})
		return
	}

	user := model.User{
		Username:     req.Username,
		PasswordHash: string(hash),
		DisplayName:  req.Username, // ��ʼ��ʾ��Ϊ�˺�
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "�����û�ʧ��"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ע��ɹ�", "id": user.ID})
}

// ��ѧ��ע�ͣ� Login ʾ����POST /login??
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "��������"})
		return
	}

	var user model.User
	if err := db.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "�û������������"})
		return
	}

	// ��ѧ��ע�ͣ� �ȶ� bcrypt ��ϣ
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "�û������������"})
		return
	}

	// ��ѧ��ע�ͣ� ��¼�ɹ� -> ���� token����ԭ��??JWT �߼��ͷ�����??
	// ��ѧ��ע�ͣ� ���� JWT
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "���� token ʧ��"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
