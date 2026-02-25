package router

import (
	"project/internal/db"
	"project/internal/handler"
	"project/internal/middleware"
	"project/internal/model"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.RequestLogger())

	if err := db.DB.AutoMigrate(
		&model.User{},
		&model.Article{},
		&model.Comment{},
		&model.Like{},
		&model.Follow{},
	); err != nil {
	}

	r.Static("/static", "./static")

	r.StaticFile("/", "./static/index.html")

	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)

	r.GET("/articles", handler.ListArticles)
	r.GET("/articles/:id", handler.GetArticle)
	r.GET("/articles/:id/comments", handler.ListComments)

	r.GET("/users/:id/profile", handler.GetUserProfile)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())

	auth.POST("/articles", handler.CreateArticle)
	auth.PUT("/articles/:id", handler.UpdateArticle)
	auth.DELETE("/articles/:id", handler.DeleteArticle)

	auth.POST("/articles/:id/comments", handler.PostComment)
	auth.POST("/articles/:id/like", handler.ToggleArticleLike)
	auth.POST("/comments/:id/like", handler.ToggleCommentLike)

	auth.POST("/users/:id/follow", handler.ToggleFollow)
	auth.GET("/me/follows", handler.GetMyFollows)

	auth.GET("/me/profile", handler.GetMyProfile)
	auth.PUT("/me/profile", handler.UpdateMyProfile)

	return r
}
