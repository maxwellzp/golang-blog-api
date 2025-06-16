package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"maxwellzp/blog-api/internal/config"
	"maxwellzp/blog-api/internal/handler"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func registerRoutes(e *echo.Echo, cfg *config.Config, log *zap.SugaredLogger, auth *handler.AuthHandler,
	blog *handler.BlogHandler,
	comment *handler.CommentHandler) {

	// Middleware
	// Recover middleware recovers from panics anywhere in the chain, prints stack trace
	e.Use(middleware.Recover())

	// Secure middleware provides protection against cross-site scripting (XSS) attack, content type sniffing,
	// clickjacking, insecure connection and other code injection attacks.
	e.Use(middleware.Secure())

	// Body limit middleware sets the maximum allowed size for a request body
	e.Use(middleware.BodyLimit(cfg.BodyLimit))

	// Custom logger middleware with zap
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

	// Health
	e.GET("/healthz", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	// Auth
	e.POST("/register", auth.Register)
	e.POST("/login", auth.Login)

	// Blogs
	e.POST("/blogs", blog.Create)
	e.GET("/blogs", blog.List)
	e.GET("/blogs/:id", blog.GetByID)
	e.PUT("/blogs/:id", blog.Update)
	e.DELETE("/blogs/:id", blog.Delete)

	// Comments
	e.POST("/comments", comment.Create)
	e.GET("/comments/:id", comment.GetByID)
	e.PUT("/comments/:id", comment.Update)
	e.DELETE("/comments/:id", comment.Delete)
	e.GET("/blogs/:blog_id/comments", comment.ListByBlogID)
}
