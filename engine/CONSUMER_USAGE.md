# Kafka Consumer Usage Guide

This guide explains how to use the Kafka consumer to read Debezium CDC events from Kafka.

## Overview

The consumer reads change data capture (CDC) events from PostgreSQL via Debezium and processes them using a callback handler pattern.

## Key Components

### Event Structure

```go
type Event struct {
    Operation string        // "c" (create), "u" (update), "d" (delete), "r" (read/snapshot)
    Before    *objects.User // State before the change (nil for creates)
    After     *objects.User // State after the change (nil for deletes)
    Source    SourceInfo    // Metadata like table name, timestamp, etc.
    Timestamp time.Time     // When the event was created
}
```

### Operation Types

- **"c"** (Create): A new user was inserted. Check `event.After` for the new user data.
- **"u"** (Update): An existing user was updated. Both `event.Before` and `event.After` contain data.
- **"d"** (Delete): A user was deleted. Check `event.Before` for the deleted user data.
- **"r"** (Read): Initial snapshot data loaded when Debezium starts. Check `event.After`.

## Basic Usage Example

```go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/ahsansaif47/blockchain-address-watcher/engine/config"
    "github.com/ahsansaif47/blockchain-address-watcher/engine/consumer"
)

func main() {
    // Create context for graceful shutdown
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Set up signal handling
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-sigChan
        log.Println("Received shutdown signal, stopping consumer...")
        cancel()
    }()

    // Load Kafka configuration
    cfg, err := config.ConsumerConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Create Kafka manager
    km, err := consumer.NewKafkaManager(&consumer.Config{
        Broker:          cfg.Broker,
        Topic:           cfg.Topic,
        Partition:       cfg.Partition,
        MaxRetries:      cfg.MaxRetries,
        RetryDelay:      cfg.RetryDelay * time.Second,
        HealthCheckFreq: cfg.HealthCheckFreq * time.Second,
    })
    if err != nil {
        log.Fatalf("Failed to create Kafka manager: %v", err)
    }
    defer km.Close()

    // Define event handler
    handleEvent := func(event *consumer.Event) error {
        switch event.Operation {
        case "c": // Create
            if event.After != nil {
                log.Printf("[CREATE] New user: %s (Wallet: %s)", 
                    event.After.Email, event.After.WalletAddress)
                // TODO: Add wallet to blockchain watching list
            }

        case "u": // Update
            if event.Before != nil && event.After != nil {
                log.Printf("[UPDATE] User %s updated", event.After.Email)
                
                // Check if wallet address changed
                if event.Before.WalletAddress != event.After.WalletAddress {
                    log.Printf("  Wallet changed: %s -> %s", 
                        event.Before.WalletAddress, event.After.WalletAddress)
                    // TODO: Update watching list
                }
            }

        case "d": // Delete
            if event.Before != nil {
                log.Printf("[DELETE] User deleted: %s", event.Before.Email)
                // TODO: Remove wallet from watching list
            }

        case "r": // Read (snapshot)
            if event.After != nil {
                log.Printf("[SNAPSHOT] Loading user: %s", event.After.Email)
                // TODO: Add existing wallets to watching list
            }
        }
        
        return nil
    }

    // Start reading from Kafka
    log.Println("Starting Kafka consumer...")
    if err := consumer.ReadWithRetry(ctx, km, handleEvent, 5*time.Second); err != nil {
        log.Printf("Consumer stopped: %v", err)
    }

    log.Println("Consumer stopped successfully")
}
```

## Advanced Handler Example

```go
// Advanced handler with business logic
func advancedHandler() consumer.EventHandler {
    return func(event *consumer.Event) error {
        // Integrate with blockchain monitoring service
        switch event.Operation {
        case "c", "u":
            if event.After != nil && event.After.Subscribed && event.After.WalletAddress != "" {
                // Only watch subscribed users with wallet addresses
                err := monitorBlockchainAddress(event.After.WalletAddress, event.After.Email)
                if err != nil {
                    return fmt.Errorf("failed to monitor wallet: %w", err)
                }
            }
            
        case "d":
            if event.Before != nil && event.Before.WalletAddress != "" {
                err := stopMonitoringBlockchainAddress(event.Before.WalletAddress)
                if err != nil {
                    return fmt.Errorf("failed to stop monitoring: %w", err)
                }
            }
        }
        
        return nil
    }
}

func monitorBlockchainAddress(address, userEmail string) error {
    // Implement actual blockchain monitoring logic
    // This could call an external API, set up webhooks, etc.
    fmt.Printf("Starting monitoring for wallet %s (user: %s)\n", address, userEmail)
    return nil
}

func stopMonitoringBlockchainAddress(address string) error {
    // Implement actual stop monitoring logic
    fmt.Printf("Stopping monitoring for wallet %s\n", address)
    return nil
}
```

## Configuration

Ensure your `.env` file contains the following variables:

```env
KAFKA_BROKER=localhost:9092
KAFKA_TOPIC=<your-debezium-topic-name>
```

The Kafka topic name is determined by your Debezium connector configuration. It typically follows the format: `<database_server_name>.<schema_name>.<table_name>`

For example, if your Debezium connector is named `postgres-connector` and you're watching the `public.users` table, the topic would be: `postgres-connector.public.users`

## Running the Consumer

To run the consumer:

```bash
cd engine
go run main.go
```

Make sure your Kafka, PostgreSQL, and Debezium Connect services are running:

```bash
docker-compose up -d kafka postgres debezium
```

## Features

- **Automatic Reconnection**: The KafkaManager handles connection failures and automatically reconnects
- **Health Checks**: Periodic health checks ensure the connection remains active
- **Graceful Shutdown**: Supports context cancellation for clean shutdown
- **Retry Logic**: ReadWithRetry provides automatic retry on connection failures
- **Error Handling**: Individual event failures don't stop the consumer
- **Offset Management**: Automatically commits offsets after successful processing

## Testing

You can test the consumer by making changes to your PostgreSQL database:

```sql
-- Create a new user (will trigger 'c' operation)
INSERT INTO users (email, password_hash, phone_no, wallet_address, subscribed)
VALUES ('test@example.com', 'hash', '1234567890', '0x123...', true);

-- Update a user (will trigger 'u' operation)
UPDATE users SET wallet_address = '0x456...' WHERE email = 'test@example.com';

-- Delete a user (will trigger 'd' operation)
DELETE FROM users WHERE email = 'test@example.com';
```

## Troubleshooting

### No messages received
- Verify Debezium is running and the connector is configured
- Check that the topic name matches your Debezium configuration
- Ensure PostgreSQL has changes to capture (try making a test insert)

### Connection errors
- Verify Kafka is running: `docker-compose ps kafka`
- Check the broker address in your `.env` file
- Review Kafka logs: `docker-compose logs kafka`

### Parse errors
- Check the Debezium connector configuration to ensure it matches your table structure
- Verify the User model matches your database schema
- Check Debezium message format in the topic: `kafka-console-consumer --bootstrap-server localhost:9092 --topic <topic-name> --from-beginning`

## Integration with Blockchain Watching

The consumer is designed to integrate with your blockchain watching system:

1. **Create/Update**: When a user is created or updated with a wallet address, add it to your watching list
2. **Delete**: When a user is deleted, remove their wallet from watching
3. **Snapshot**: Load existing users' wallets when the consumer first starts

You can implement the actual blockchain monitoring logic in your handler function, calling external APIs or services as needed.
