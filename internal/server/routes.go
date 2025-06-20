package server

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
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
	loginLimiter := echoMiddleware.NewRateLimiterMemoryStoreWithConfig(
		echoMiddleware.RateLimiterMemoryStoreConfig{
			Rate:      rate.Every(30 * time.Second), // 1 request every 30s
			Burst:     2,                            // How many tokens can be used immediately in a short burst
			ExpiresIn: 15 * time.Minute,             // Sets how long the client’s token bucket state is remembered in memory (per IP by default).
		},
	)

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
	e.POST("/login", auth.Login, echoMiddleware.RateLimiter(loginLimiter))
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
