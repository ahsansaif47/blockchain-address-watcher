package main

import (
	"log"
	"os"
	"time"

	"github.com/ahsansaif47/blockchain-address-watcher/engine/consumer"
	"github.com/ahsansaif47/blockchain-address-watcher/engine/utils"
	"github.com/joho/godotenv"
)

func consumerConfig() (*consumer.Config, error) {
	if err := godotenv.Load("..", ".env"); err != nil {
		return nil, err
	}

	partition, err := utils.StringToInteger(os.Getenv("KAFKA_PARTITION"))
	if err != nil {
		return nil, err
	}
	retries, err := utils.StringToInteger(os.Getenv("KAFKA_RETRIES"))
	if err != nil {
		return nil, err
	}
	delay, err := utils.StringToInteger(os.Getenv("KAFKA_RETRY_DELAY"))
	if err != nil {
		return nil, err
	}
	healthFreq, err := utils.StringToInteger(os.Getenv("KAFKA_HEALTH_FREQ"))
	if err != nil {
		return nil, err
	}
	return &consumer.Config{
		Broker:          os.Getenv("KAFKA_BROKER"),
		Topic:           os.Getenv("KAFKA_TOPIC"),
		Partition:       partition,
		MaxRetries:      retries,
		RetryDelay:      time.Duration(delay),
		HealthCheckFreq: time.Duration(healthFreq),
	}, nil
}

func main() {

	config, err := consumerConfig()
	if err != nil {
		log.Fatalf("Error getting kafka manager config: %v", err)
	}

	km, err := consumer.NewKafkaManager(config)
	if err != nil {
		log.Fatalf("Error creating kafka manager: %v", err)
	}
	defer km.Close()
}
