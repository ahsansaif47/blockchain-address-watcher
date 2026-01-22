package main

import (
	"log"

	"github.com/ahsansaif47/blockchain-address-watcher/engine/config"
	"github.com/ahsansaif47/blockchain-address-watcher/engine/consumer"
)

func main() {
	config, err := config.ConsumerConfig()
	if err != nil {
		log.Fatalf("Error getting kafka manager config: %v", err)
	}

	km, err := consumer.NewKafkaManager(config)
	if err != nil {
		log.Fatalf("Error creating kafka manager: %v", err)
	}
	defer km.Close()
}
