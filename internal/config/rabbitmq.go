package config

import (
	"errors"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
)

func RabbitMQConnection() (*amqp091.Connection, error) {
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file not found, using system environment variables")
	}

	mqUrl := os.Getenv("MQ_HOST")
	if mqUrl == "" {
		log.Fatalf("MQ_HOST environment variable is not set")
		return nil, errors.New("MQ_HOST environment variable is not set")
	}

	conn, err := amqp091.Dial(mqUrl)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		return nil, err
	}
	log.Println("Connected to RabbitMQ successfully")

	return conn, nil
}
