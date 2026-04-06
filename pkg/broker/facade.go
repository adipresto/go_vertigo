package broker

import (
	"database/sql"
	"log"
	"vertigo/pkg/config"
	"vertigo/pkg/db"
	"vertigo/pkg/network"
)

// TripleBaseBroker is our Master Facade (Base 1: DB, Base 2: Centrifugo, Base 3: MQTT)
type TripleBaseBroker struct {
	DB   *sql.DB
	Net  network.Publisher
	MQTT network.Publisher
}

// NewBroker initializes the Facade, hiding the complexity of subsystems
func NewBroker(cfg *config.Config) (*TripleBaseBroker, error) {
	pool, err := db.NewPool(cfg.Database.Path)
	if err != nil {
		return nil, err
	}

	broker := &TripleBaseBroker{DB: pool}

	// Initialize Base 2 (Centrifugo) if enabled
	if cfg.Network.Centrifugo.Enabled {
		messenger, err := network.NewMessenger(cfg.Network.Centrifugo.URL)
		if err != nil {
			log.Printf("Warning: Base 2 (Centrifugo) failed to connect: %v", err)
		} else {
			broker.Net = messenger
		}
	}

	// Initialize Base 3 (MQTT) if enabled
	if cfg.Network.MQTT.Enabled {
		mqttPub, err := network.NewMQTTPublisher(cfg.Network.MQTT.URL)
		if err != nil {
			log.Printf("Warning: Base 3 (MQTT) failed to connect: %v", err)
		} else {
			broker.MQTT = mqttPub
		}
	}

	return broker, nil
}
