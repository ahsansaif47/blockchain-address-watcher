package parser

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ahsansaif47/blockchain-address-watcher/engine/consumer"
)

// parseDebeziumMessage parses a raw Debezium message into an Event struct
func parseDebeziumMessage(data []byte) (*consumer.Event, error) {
	var msg consumer.DebeziumMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Debezium message: %w", err)
	}

	// Validate operation type
	operation := msg.Payload.Operation
	if operation == "" {
		return nil, fmt.Errorf("missing operation type in payload")
	}

	// Create event
	event := &consumer.Event{
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
