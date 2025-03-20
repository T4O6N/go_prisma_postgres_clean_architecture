package config

import (
	"fmt"
	"log/slog"
	"os"
	"sample-project/prisma/db"

	"github.com/joho/godotenv"
)

func ConnectDB() (*db.PrismaClient, error) {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file not found, using system environment variables")
	}

	// Check if DATABASE_URL is set
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		slog.Error("DATABASE_URL environment variable is not set")
		return nil, ErrMissingDatabaseURL
	}

	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		slog.Error("Failed to connect to database", "error", err)
		return nil, err
	}

	slog.Info("Database connected successfully")
	return client, nil

}

var ErrMissingDatabaseURL = fmt.Errorf("DATABASE_URL is required but not set")
