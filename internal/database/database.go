package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"maxwellzp/blog-api/internal/config"
)

func Connect(cfg *config.Config, logger *zap.SugaredLogger) *sql.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.MySQLUser,
		cfg.MySQLPassword,
		cfg.MySQLHost,
		cfg.MySQLPort,
		cfg.MySQLDatabase,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Fatalw("failed to connect to database",
			"error", err,
			"dsn", dsn,
		)
	}

	if err = db.Ping(); err != nil {
		logger.Fatalw("failed to ping database",
			"error", err,
		)
	}

	logger.Infow("successfully connected to database")
	return db
}
