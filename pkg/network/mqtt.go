package network

import (
	"context"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQTTPublisher is the MQTT implementation of Publisher
type MQTTPublisher struct {
	Client mqtt.Client
}

func NewMQTTPublisher(brokerURL string) (*MQTTPublisher, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerURL)
	opts.SetClientID("vertigo-publisher-" + fmt.Sprintf("%d", time.Now().Unix()))
	opts.SetAutoReconnect(true)

	opts.OnConnect = func(c mqtt.Client) {
		log.Printf("Connected to MQTT Broker: %s", brokerURL)
	}
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Printf("Lost connection to MQTT Broker: %v", err)
	}

	client := mqtt.NewClient(opts)
	token := client.Connect()
	
	// We use a short timeout for the initial connection attempt
	if token.WaitTimeout(2 * time.Second) && token.Error() != nil {
		return nil, fmt.Errorf("failed to connect to MQTT: %v", token.Error())
	}

	return &MQTTPublisher{Client: client}, nil
}

func (m *MQTTPublisher) Publish(ctx context.Context, topic string, data []byte) error {
	if !m.Client.IsConnected() {
		return fmt.Errorf("mqtt client not connected")
	}
	
	token := m.Client.Publish(topic, 1, false, data)
	if token.WaitTimeout(1 * time.Second) && token.Error() != nil {
		return token.Error()
	}
	
	return nil
}
