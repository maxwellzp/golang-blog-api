package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"maxwellzp/blog-api/internal/config"
	"maxwellzp/blog-api/internal/database"
	"maxwellzp/blog-api/internal/handler"
	"maxwellzp/blog-api/internal/repository"
	"maxwellzp/blog-api/internal/service"
	"maxwellzp/blog-api/internal/validation"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	e    *echo.Echo
	cfg  *config.Config
	log  *zap.SugaredLogger
	port string
}

func New(cfg *config.Config, logger *zap.SugaredLogger) *Server {
	e := echo.New()

	db := database.Connect(cfg, logger)
	logger.Infow("connected to database",
		"host", cfg.MySQLHost,
		"port", cfg.MySQLPort,
	)
	validator := validation.NewValidator()

	// DI
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService, logger, validator)

	blogRepo := repository.NewBlogRepository(db)
	blogService := service.NewBlogService(blogRepo)
	blogHandler := handler.NewBlogHandler(blogService, logger, validator)

	commentRepo := repository.NewCommentRepository(db)
	commentService := service.NewCommentService(commentRepo)
	commentHandler := handler.NewCommentHandler(commentService, logger, validator)

	// Routes + Middleware
	registerRoutes(e, cfg, logger, authHandler, blogHandler, commentHandler)

	return &Server{
		e:    e,
		cfg:  cfg,
		log:  logger,
		port: cfg.ServerPort,
	}
}

func (s *Server) Start() {
	s.log.Infow("starting server",
		"port", s.port,
	)

	// - SIGINT (interrupt signal, like Ctrl+C)
	// - SIGTERM (termination signal, like docker stop)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := s.e.Start(":" + s.port); err != nil && err != http.ErrServerClosed {
			s.log.Errorw("server error",
				"error", err,
			)
		}
	}()

	<-ctx.Done()
	s.log.Infow("shutting down server...")

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
	if err := s.e.Shutdown(shutdownCtx); err != nil {
		s.log.Errorw("error occurred on server shutdown",
			"error", err,
		)
	}
}
