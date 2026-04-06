package network

import (
	"context"
	"fmt"
	"log"

	"github.com/centrifugal/centrifuge-go"
)

// Publisher defines the interface for real-time messaging
type Publisher interface {
	Publish(ctx context.Context, channel string, data []byte) error
}

// Messenger is the Centrifugo implementation of Publisher
type Messenger struct {
	Client *centrifuge.Client
}

func NewMessenger(endpoint string) (*Messenger, error) {
	c := centrifuge.NewJsonClient(endpoint, centrifuge.Config{})

	c.OnConnecting(func(e centrifuge.ConnectingEvent) {
		log.Printf("Connecting to Centrifugo: %s", endpoint)
	})

	c.OnConnected(func(e centrifuge.ConnectedEvent) {
		log.Printf("Connected to Centrifugo as %s", e.ClientID)
	})

	c.OnDisconnected(func(e centrifuge.DisconnectedEvent) {
		log.Printf("Disconnected from Centrifugo: %s", e.Reason)
	})

	err := c.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Centrifugo: %v", err)
	}

	return &Messenger{Client: c}, nil
}

func (m *Messenger) Publish(ctx context.Context, channel string, data []byte) error {
	if m.Client == nil {
		return fmt.Errorf("centrifuge client is nil")
	}
	_, err := m.Client.Publish(ctx, channel, data)
	return err
}
