package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	objects "github.com/ahsansaif47/blockchain-address-watcher/engine/models"
	"github.com/segmentio/kafka-go"
)

// Event represents a parsed Debezium CDC event
type Event struct {
	Operation string        // "c" (create), "u" (update), "d" (delete), "r" (read/snapshot)
	Before    *objects.User // State before the change (nil for creates)
	After     *objects.User // State after the change (nil for deletes)
	Source    SourceInfo    // Metadata like table name, timestamp, etc.
	Timestamp time.Time     // When the event was created
}

// SourceInfo contains metadata from Debezium about the source of the event
type SourceInfo struct {
	Version   string `json:"version"`
	Connector string `json:"connector"`
	Name      string `json:"name"`
	TsMs      int64  `json:"ts_ms"`
	TsUs      int64  `json:"ts_us"`
	TsNs      int64  `json:"ts_ns"`
	Snapshot  bool   `json:"snapshot"` // "true", "false", or "last"
	Db        string `json:"db"`
	Sequence  string `json:"sequences"`
	Schema    string `json:"schema"`
	Table     string `json:"table"`
	TxId      string `json:"txId,omitempty"`
	Lsn       int64  `json:"lsn,omitempty"`
}

// DebeziumMessage represents the raw Debezium message structure
type DebeziumMessage struct {
	Schema  DebeziumSchema  `json:"schema"`
	Payload DebeziumPayload `json:"payload"`
}

// DebeziumSchema contains the schema information from Debezium
type DebeziumSchema struct {
	Type     string                `json:"type"`
	Fields   []DebeziumSchemaField `json:"fields"`
	Optional bool                  `json:"optional"`
	Name     string                `json:"name"`
	Version  int                   `json:"version"`
}

// DebeziumSchemaField represents a field in the schema
type DebeziumSchemaField struct {
	Type     string `json:"type"`
	Field    string `json:"field"`
	Optional bool   `json:"optional"`
}

// DebeziumPayload contains the actual data from Debezium
type DebeziumPayload struct {
	Before    *objects.User `json:"before"`
	After     *objects.User `json:"after"`
	Source    SourceInfo    `json:"source"`
	Operation string        `json:"op"`
	TsMs      int64         `json:"ts_ms"`
	TsUs      int64         `json:"ts_us"`
	TsNs      int64         `json:"ts_ns"`
}

// EventHandler is a callback function that processes each Debezium event
// It receives the parsed event and returns an error if processing fails
type EventHandler func(event *Event) error

// Read continuously consumes messages from Kafka and processes them using the provided handler
// Parameters:
//   - ctx: Context for cancellation and graceful shutdown
//   - km: KafkaManager instance for connection management
//   - handler: Callback function to process each event
//
// Example usage:
//
//	func handleEvent(event *consumer.Event) error {
//	    switch event.Operation {
//	    case "c", "u":
//	        fmt.Printf("User created/updated: %s\n", event.After.Email)
//	    case "d":
//	        fmt.Printf("User deleted: %s\n", event.Before.Email)
//	    }
//	    return nil
//	}
//
//	ctx := context.Background()
//	consumer.Read(ctx, kafkaManager, handleEvent)
func Read(ctx context.Context, km *KafkaManager, handler EventHandler) error {
	if km == nil {
		return fmt.Errorf("KafkaManager cannot be nil")
	}
	if handler == nil {
		return fmt.Errorf("event handler cannot be nil")
	}

	// Create a reader for the topic
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{km.config.Broker},
		Topic:    km.config.Topic,
		GroupID:  "blockchain-address-watcher-group",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	defer r.Close()

	log.Printf("[Reader] Starting to read from topic: %s", km.config.Topic)

	// Start reading loop
	for {
		select {
		case <-ctx.Done():
			log.Printf("[Reader] Context cancelled, stopping reader")
			return ctx.Err()

		default:
			// Read message from Kafka
			m, err := r.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					log.Printf("[Reader] Context cancelled during read: %v", err)
					return ctx.Err()
				}
				log.Printf("[Reader] Error reading message: %v", err)
				// Continue trying to read
				time.Sleep(1 * time.Second)
				continue
			}

			log.Printf("[Reader] Received message at offset %d (partition %d)",
				m.Offset, m.Partition)

			// Parse the Debezium message
			event, err := parseDebeziumMessage(m.Value)
			if err != nil {
				log.Printf("[Reader] Error parsing message: %v", err)
				continue
			}

			// Call the event handler
			if err := handler(event); err != nil {
				log.Printf("[Reader] Error in event handler: %v", err)
				// Continue processing other messages even if one fails
			}
		}
	}
}

// parseDebeziumMessage parses a raw Debezium message into an Event struct
func parseDebeziumMessage(data []byte) (*Event, error) {
	var msg DebeziumMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Debezium message: %w", err)
	}

	// Validate operation type
	operation := msg.Payload.Operation
	if operation == "" {
		return nil, fmt.Errorf("missing operation type in payload")
	}

	// Create event
	event := &Event{
		Operation: operation,
		Before:    msg.Payload.Before,
		After:     msg.Payload.After,
		Source:    msg.Payload.Source,
		Timestamp: time.UnixMilli(msg.Payload.TsMs),
	}

	// Validate event data
	switch operation {
	case "c", "r": // Create or Read (snapshot)
		if event.After == nil {
			return nil, fmt.Errorf("missing 'after' data for operation '%s'", operation)
		}
	case "u": // Update
		if event.Before == nil || event.After == nil {
			return nil, fmt.Errorf("missing 'before' or 'after' data for operation 'u'")
		}
	case "d": // Delete
		if event.Before == nil {
			return nil, fmt.Errorf("missing 'before' data for operation 'd'")
		}
	default:
		return nil, fmt.Errorf("unknown operation type: %s", operation)
	}

	return event, nil
}

// ReadWithRetry wraps the Read function with automatic retry logic
// It will retry reading indefinitely if the connection is lost
func ReadWithRetry(ctx context.Context, km *KafkaManager, handler EventHandler, retryDelay time.Duration) error {
	if retryDelay == 0 {
		retryDelay = 5 * time.Second
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := Read(ctx, km, handler)
			if err != nil {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				log.Printf("[ReaderWithRetry] Read failed: %v, retrying in %v", err, retryDelay)
				time.Sleep(retryDelay)
				continue
			}
			return nil
		}
	}
}
