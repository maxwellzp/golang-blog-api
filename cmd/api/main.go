package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"maxwellzp/blog-api/internal/config"
	"maxwellzp/blog-api/internal/database"
	"maxwellzp/blog-api/internal/handler"
	"maxwellzp/blog-api/internal/logger"
	"maxwellzp/blog-api/internal/repository"
	"maxwellzp/blog-api/internal/service"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logr, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logr.Sync()

	cfg := config.Load()
	logr.Infow("starting server...",
		"port", cfg.ServerPort,
	)

	db := database.Connect(cfg)
	defer db.Close()
	logr.Infow("connected to database",
		"db_host", cfg.MySQLHost,
		"db_port", cfg.MySQLPort,
	)

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService, logr)

	blogRepo := repository.NewBlogRepository(db)
	blogService := service.NewBlogService(blogRepo)
	blogHandler := handler.NewBlogHandler(blogService, logr)

	commentRepo := repository.NewCommentRepository(db)
	commentService := service.NewCommentService(commentRepo)
	commentHandler := handler.NewCommentHandler(commentService, logr)

	e := echo.New()

	// Middleware
	// Recover middleware recovers from panics anywhere in the chain, prints stack trace
	e.Use(middleware.Recover())

	// Secure middleware provides protection against cross-site scripting (XSS) attack, content type sniffing,
	// clickjacking, insecure connection and other code injection attacks.
	e.Use(middleware.Secure())

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			logr.Infow("request",
				"method", c.Request().Method,
				"path", c.Request().URL.Path,
				"status", c.Response().Status,
				"latency", time.Since(start),
				"user_agent", c.Request().UserAgent(),
			)
			return err
		}
	})

	// Body limit middleware sets the maximum allowed size for a request body
	e.Use(middleware.BodyLimit(cfg.BodyLimit))

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

	e.GET("/healthz", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Runs the Echo server in a separate goroutine so the main thread can continue.
	go func() {
		if err := e.Start(":" + cfg.ServerPort); err != nil && err != http.ErrServerClosed {
			logr.Errorw("shutting down the server",
				"error", err,
			)
		}
	}()

	//quit := make(chan os.Signal, 1)
	//signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// signal.Notify(...) makes the program listen for:
	// - SIGINT (interrupt signal, like Ctrl+C)
	// - SIGTERM (termination signal, like docker stop)
	// <-quit // blocks until one of these signals is received.

	// Wait for signal
	<-ctx.Done()
	logr.Infow("shutting down server...")

	/*
		Creates a context with a 10-second timeout — this gives the server time to clean up
		(close DB connections, finish ongoing requests).
		Avoids killing in-progress HTTP requests.
		Prevents corrupted states (e.g., half-written DB rows).
		Works well with containers, Kubernetes, systemd, etc.
	*/
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Calls e.Shutdown(shutdownCtx) to gracefully stop the Echo server.
	if err := e.Shutdown(shutdownCtx); err != nil {
		logr.Errorw("error occurred on server shutdown",
			"error", err,
		)
	}
}
