package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"maxwellzp/blog-api/internal/config"
	"maxwellzp/blog-api/internal/database"
	"maxwellzp/blog-api/internal/handler"
	"maxwellzp/blog-api/internal/repository"
	"maxwellzp/blog-api/internal/service"
)

func main() {
	cfg := config.Load()

	db := database.Connect(cfg)
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	blogRepo := repository.NewBlogRepository(db)
	blogService := service.NewBlogService(blogRepo)
	blogHandler := handler.NewBlogHandler(blogService)

	commentRepo := repository.NewCommentRepository(db)
	commentService := service.NewCommentService(commentRepo)
	commentHandler := handler.NewCommentHandler(commentService)

	e := echo.New()

	// Middleware
	// Recover middleware recovers from panics anywhere in the chain, prints stack trace
	e.Use(middleware.Recover())

	// Secure middleware provides protection against cross-site scripting (XSS) attack, content type sniffing,
	// clickjacking, insecure connection and other code injection attacks.
	e.Use(middleware.Secure())

	// Logger middleware logs the information about each HTTP request.
	e.Use(middleware.Logger())

	// Body limit middleware sets the maximum allowed size for a request body
	e.Use(middleware.BodyLimit("1M"))

	e.POST("/register", authHandler.Register)
	e.POST("/login", authHandler.Login)

	e.POST("/blogs", blogHandler.Create)
	e.GET("/blogs", blogHandler.List)
	e.GET("/blogs/:id", blogHandler.GetByID)
	e.PUT("/blogs/:id", blogHandler.Update)
	e.DELETE("/blogs/:id", blogHandler.Delete)

	e.POST("/comments", commentHandler.Create)
	e.GET("/comments/:id", commentHandler.GetByID)
	e.PUT("/comments/:id", commentHandler.Update)
	e.DELETE("/comments/:id", commentHandler.Delete)
	e.GET("/blogs/:blog_id/comments", commentHandler.ListByBlogID)

	e.Start(":" + cfg.ServerPort)
}
