package config

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
)

type Config struct {
	ServerPort    string
	MySQLUser     string
	MySQLPassword string
	MySQLHost     string
	MySQLPort     string
	MySQLDatabase string
	JWTSecret     string
	BodyLimit     string
}

func Load(logger *zap.SugaredLogger) *Config {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logger.Warnw("No .env file found")
	}
	return &Config{
		ServerPort:    getEnv(logger, "SERVER_PORT", "8080"),
		MySQLUser:     mustGetEnv(logger, "MYSQL_USER"),
		MySQLPassword: mustGetEnv(logger, "MYSQL_PASSWORD"),
		MySQLHost:     mustGetEnv(logger, "MYSQL_HOST"),
		MySQLPort:     getEnv(logger, "MYSQL_PORT", "3306"),
		MySQLDatabase: mustGetEnv(logger, "MYSQL_DATABASE"),
		JWTSecret:     mustGetEnv(logger, "JWT_SECRET"),
		BodyLimit:     getEnv(logger, "BODY_LIMIT", "1M"),
	}
}

func getEnv(logger *zap.SugaredLogger, key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	logger.Infow("using default value for env variable",
		"key", key,
		"default", defaultVal,
	)
	return defaultVal
}

// mustGetEnv - Distinguishes between "unset" and "empty"
func mustGetEnv(logger *zap.SugaredLogger, key string) string {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		logger.Fatalw("required environment variable is not set",
			"key", key,
		)
	}
	return val
}
