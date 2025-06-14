package config

import (
	"github.com/joho/godotenv"
	"log"
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
}

func Load() *Config {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	return &Config{
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		MySQLUser:     mustGetEnv("MYSQL_USER"),
		MySQLPassword: mustGetEnv("MYSQL_PASSWORD"),
		MySQLHost:     mustGetEnv("MYSQL_HOST"),
		MySQLPort:     getEnv("MYSQL_PORT", "3306"),
		MySQLDatabase: mustGetEnv("MYSQL_DATABASE"),
		JWTSecret:     mustGetEnv("JWT_SECRET"),
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

// mustGetEnv - Distinguishes between "unset" and "empty"
func mustGetEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return val
}
