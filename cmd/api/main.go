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

	e := echo.New()
	e.POST("/register", authHandler.Register)
	e.POST("/login", authHandler.Login)
}
