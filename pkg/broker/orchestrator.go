package broker

import (
	"bytes"
	"context"
	"fmt"
	"vertigo/pkg/db"
	"vertigo/pkg/model"
)

// Dispatch is the single entry point for the business logic (SQL Abstraction)
func (b *TripleBaseBroker) Dispatch(ctx context.Context, sql string, channel string, args ...any) ([]byte, error) {
	// 1. Stream from Database (Base 1)
	var buf bytes.Buffer
	weight, err := db.StreamingQuery(ctx, b.DB, sql, &buf, args...)
	if err != nil {
		return nil, fmt.Errorf("DB Error (Base 1): %v", err)
	}

	// 2. Wrap into Final Payload (DTO)
	payload := model.WrapPayload(sql, weight, buf.Bytes())
	
	// 3. Broadcast to Network (Base 2: Centrifugo)
	if b.Net != nil {
		go func() {
			err := b.Net.Publish(context.Background(), channel, payload)
			if err != nil {
				fmt.Printf("Base 2 (Centrifugo) Broadcast Error: %v\n", err)
			}
		}()
	}

	// 4. Broadcast to Network (Base 3: MQTT)
	if b.MQTT != nil {
		go func() {
			// Using same channel name as MQTT Topic
			err := b.MQTT.Publish(context.Background(), channel, payload)
			if err != nil {
				fmt.Printf("Base 3 (MQTT) Broadcast Error: %v\n", err)
			}
		}()
	}

	return payload, nil
}
