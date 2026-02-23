package router

import (
	"project/internal/db"
	"project/internal/handler"
	"project/internal/middleware"
	"project/internal/model"

	"github.com/gin-gonic/gin"
)

// SetupRouter 创建并返回 gin Engine，包含静态资源和所有路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 自动建表（若表已存在则不会覆盖）
	if err := db.DB.AutoMigrate(
		&model.User{},
		&model.Article{},
		&model.Comment{},
		&model.Like{},
		&model.Follow{},
	); err != nil {
		// 可记录错误
	}

	// 静态文件（前端资源）
	r.Static("/static", "./static")

	// 公共（无需鉴权）路由
	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)

	// 文章相关（公开查询）
	r.GET("/articles", handler.ListArticles)
	r.GET("/articles/:id", handler.GetArticle)
	r.GET("/articles/:id/comments", handler.ListComments)

	// 公开查看用户资料
	r.GET("/users/:id/profile", handler.GetUserProfile)

	// 需要鉴权的接口组
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())

	// 文章增删改（需要登录）
	auth.POST("/articles", handler.CreateArticle)
	auth.PUT("/articles/:id", handler.UpdateArticle)
	auth.DELETE("/articles/:id", handler.DeleteArticle)

	// 评论、点赞、关注等需要登录的操作
	auth.POST("/articles/:id/comments", handler.PostComment)
	auth.POST("/articles/:id/like", handler.ToggleArticleLike)
	auth.POST("/comments/:id/like", handler.ToggleCommentLike)

	// 关注相关
	auth.POST("/users/:id/follow", handler.ToggleFollow)
	auth.GET("/me/follows", handler.GetMyFollows)

	// 当前用户 profile 获取/更新
	auth.GET("/me/profile", handler.GetMyProfile)
	auth.PUT("/me/profile", handler.UpdateMyProfile)

	return r
}
