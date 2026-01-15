package config

import (
	"os"
	"time"

	"github.com/ahsansaif47/blockchain-address-watcher/engine/consumer"
	"github.com/joho/godotenv"
)

func ConsumerConfig() (*consumer.Config, error) {
	if err := godotenv.Load("..", ".env"); err != nil {
		return nil, err
	}

	return &consumer.Config{
		Broker:          os.Getenv("KAFKA_BROKER"),
		Topic:           os.Getenv("KAFKA_TOPIC"),
		Partition:       0,
		MaxRetries:      5,
		RetryDelay:      time.Duration(2),
		HealthCheckFreq: time.Duration(1),
	}, nil
}
