package main

import (
	"log"
	"maxwellzp/blog-api/internal/config"
	"maxwellzp/blog-api/internal/logger"
	"maxwellzp/blog-api/internal/server"
)

func main() {
	logr, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logr.Sync()

	cfg := config.Load(logr)

	srv := server.New(cfg, logr)
	srv.Start()
}
