package server

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"maxwellzp/blog-api/internal/config"
	"maxwellzp/blog-api/internal/handler"
	appMiddleware "maxwellzp/blog-api/internal/middleware"
	"net/http"
	"time"
)

func registerRoutes(
	e *echo.Echo,
	cfg *config.Config,
	log *zap.SugaredLogger,
	auth *handler.AuthHandler,
	blog *handler.BlogHandler,
	comment *handler.CommentHandler,
) {
	// Global middleware
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.Secure())
	e.Use(echoMiddleware.BodyLimit(cfg.BodyLimit))

	// Custom logging middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			log.Infow("request",
				"method", c.Request().Method,
				"path", c.Request().URL.Path,
				"status", c.Response().Status,
				"latency", time.Since(start),
				"user_agent", c.Request().UserAgent(),
			)
			return err
		}
	})

	// --- Public Routes ---
	e.GET("/healthz", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	e.POST("/register", auth.Register)
	e.POST("/login", auth.Login)
	e.GET("/blogs", blog.List)
	e.GET("/blogs/:id", blog.GetByID)
	e.GET("/blogs/:blog_id/comments", comment.ListByBlogID)
	e.GET("/comments/:id", comment.GetByID)

	// --- Protected Routes ---
	authorized := e.Group("")
	authorized.Use(appMiddleware.JWTMiddleware(cfg.JWTSecret, log))

	// Blogs (auth required)
	authorized.POST("/blogs", blog.Create)
	authorized.PUT("/blogs/:id", blog.Update)
	authorized.DELETE("/blogs/:id", blog.Delete)

	// Comments (auth required)
	authorized.POST("/comments", comment.Create)
	authorized.PUT("/comments/:id", comment.Update)
	authorized.DELETE("/comments/:id", comment.Delete)
}
