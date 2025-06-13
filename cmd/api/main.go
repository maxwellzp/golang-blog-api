package main

import (
	"github.com/labstack/echo/v4"
	"log"
	"maxwellzp/blog-api/internal/handler"
	"maxwellzp/blog-api/internal/repository"
	"maxwellzp/blog-api/internal/service"
)

import (
	"github.com/joho/godotenv"
)

func main() {

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	userRepo := repository.NewUserRepository(nil)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	blogRepo := repository.NewBlogRepository(nil)
	blogService := service.NewBlogService(blogRepo)
	blogHandler := handler.NewBlogHandler(blogService)

	e := echo.New()
	e.POST("/register", authHandler.Register)
	e.POST("/login", authHandler.Login)

	e.POST("/blogs", blogHandler.Create)
	e.GET("/blogs", blogHandler.List)
	e.GET("/blogs/:id", blogHandler.GetByID)
	e.PUT("/blogs/:id", blogHandler.Update)
	e.DELETE("/blogs/:id", blogHandler.Delete)
	
}
