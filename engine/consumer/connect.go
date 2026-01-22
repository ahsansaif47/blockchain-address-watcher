package consumer

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

// Config holds Kafka connection configuration
type Config struct {
	Broker    string
	Topic     string
	Partition int
	// Optional: TLS and SASL configuration can be added here
	MaxRetries      int
	RetryDelay      time.Duration
	HealthCheckFreq time.Duration
}

// KafkaManager manages Kafka connections with reconnection logic, health checks, and observability
type KafkaManager struct {
	conn        *kafka.Conn
	config      *Config
	mu          sync.RWMutex
	isClosed    bool
	retryCount  int
	lastConnect time.Time
	healthCheck chan struct{}
}

// NewKafkaManager creates a new Kafka connection manager with the given configuration
func NewKafkaManager(config *Config) (*KafkaManager, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Set default values
	if config.MaxRetries == 0 {
		config.MaxRetries = 5
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = 1 * time.Second
	}
	if config.HealthCheckFreq == 0 {
		config.HealthCheckFreq = 30 * time.Second
	}

	km := &KafkaManager{
		config:      config,
		isClosed:    false,
		healthCheck: make(chan struct{}),
	}

	if err := km.connectWithRetry(); err != nil {
		return nil, err
	}

	// Start health check goroutine
	go km.runHealthCheck()

	return km, nil
}

// connect establishes a new Kafka connection
func (km *KafkaManager) connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := kafka.DialLeader(ctx, "tcp", km.config.Broker, km.config.Topic, km.config.Partition)
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka: %w", err)
	}

	km.mu.Lock()
	km.conn = conn
	km.lastConnect = time.Now()
	km.mu.Unlock()

	log.Printf("[KafkaManager] Connected to %s, topic: %s, partition: %d", 
		km.config.Broker, km.config.Topic, km.config.Partition)
	
	return nil
}

// connectWithRetry attempts to connect with exponential backoff
func (km *KafkaManager) connectWithRetry() error {
	var lastErr error
	
	for i := 0; i < km.config.MaxRetries; i++ {
		if err := km.connect(); err != nil {
			lastErr = err
			km.retryCount++
			
			// Exponential backoff: delay * 2^attempt
			backoff := km.config.RetryDelay * time.Duration(1<<uint(i))
			log.Printf("[KafkaManager] Connection attempt %d/%d failed: %v, retrying in %v", 
				i+1, km.config.MaxRetries, err, backoff)
			
			time.Sleep(backoff)
			continue
		}
		return nil
	}

	return fmt.Errorf("failed to connect after %d attempts: %w", km.config.MaxRetries, lastErr)
}

// GetConnection returns the Kafka connection with automatic reconnection
func (km *KafkaManager) GetConnection() (*kafka.Conn, error) {
	km.mu.RLock()
	if km.isClosed {
		km.mu.RUnlock()
		return nil, fmt.Errorf("connection manager is closed")
	}
	conn := km.conn
	km.mu.RUnlock()

	// Lightweight health check before returning connection
	if conn != nil {
		if km.isConnectionAlive() {
			return conn, nil
		}
		log.Printf("[KafkaManager] Connection appears dead, attempting reconnection")
	}

	// Reconnection attempt
	if err := km.connectWithRetry(); err != nil {
		return nil, fmt.Errorf("failed to reconnect to Kafka: %w", err)
	}

	km.mu.RLock()
	defer km.mu.RUnlock()
	return km.conn, nil
}

// isConnectionAlive performs a lightweight check to see if connection is still valid
func (km *KafkaManager) isConnectionAlive() bool {
	km.mu.RLock()
	defer km.mu.RUnlock()

	if km.conn == nil || km.isClosed {
		return false
	}

	// Try to read broker metadata - this is lightweight
	if _, err := km.conn.Brokers(); err != nil {
		return false
	}

	return true
}

// Close closes the Kafka connection and stops health checks
func (km *KafkaManager) Close() error {
	km.mu.Lock()
	defer km.mu.Unlock()

	if km.isClosed {
		return nil
	}

	km.isClosed = true
	close(km.healthCheck)
	
	if km.conn != nil {
		log.Printf("[KafkaManager] Closing connection to %s", km.config.Broker)
		return km.conn.Close()
	}
	
	return nil
}

// HealthCheck performs a health check on the connection
func (km *KafkaManager) HealthCheck(ctx context.Context) error {
	km.mu.RLock()
	if km.isClosed {
		km.mu.RUnlock()
		return fmt.Errorf("connection manager is closed")
	}
	km.mu.RUnlock()

	if !km.isConnectionAlive() {
		return fmt.Errorf("connection is not alive")
	}

	return nil
}

// GetStats returns statistics about the connection
func (km *KafkaManager) GetStats() map[string]interface{} {
	km.mu.RLock()
	defer km.mu.RUnlock()

	stats := map[string]interface{}{
		"broker":        km.config.Broker,
		"topic":         km.config.Topic,
		"partition":     km.config.Partition,
		"is_closed":     km.isClosed,
		"retry_count":   km.retryCount,
		"last_connect":  km.lastConnect,
		"is_connected":  km.conn != nil,
	}

	if !km.lastConnect.IsZero() {
		stats["uptime_seconds"] = time.Since(km.lastConnect).Seconds()
	}

	return stats
}

// runHealthCheck periodically checks the connection health
func (km *KafkaManager) runHealthCheck() {
	ticker := time.NewTicker(km.config.HealthCheckFreq)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if km.isConnectionAlive() {
				continue
			}
			log.Printf("[KafkaManager] Health check failed, attempting reconnection")
			if err := km.connectWithRetry(); err != nil {
				log.Printf("[KafkaManager] Auto-reconnection failed: %v", err)
			}
		case <-km.healthCheck:
			return
		}
	}
}

// Connect creates a new KafkaManager instance with default configuration
// This is kept for backward compatibility but returns error instead of panicking
func Connect() (*KafkaManager, error) {
	config := &Config{
		Broker:    "localhost:9092",
		Topic:     "topic",
		Partition: 0,
	}
	
	km, err := NewKafkaManager(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka manager: %w", err)
	}
	
	return km, nil
}
