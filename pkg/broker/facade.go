package broker

import (
	"database/sql"
	"vertigo/pkg/db"
	"vertigo/pkg/network"
)

// DoubleBaseBroker is our Master Facade
type DoubleBaseBroker struct {
	DB  *sql.DB
	Net network.Publisher
}

// NewBroker initializes the Facade, hiding the complexity of subsystems
func NewBroker(dbConn string, netEndpoint string) (*DoubleBaseBroker, error) {
	pool, err := db.NewPool(dbConn)
	if err != nil {
		return nil, err
	}

	messenger, err := network.NewMessenger(netEndpoint)
	if err != nil {
		// Resilience: Base 1 is still functional even if Base 2 fails to connect
		return &DoubleBaseBroker{DB: pool, Net: nil}, err
	}

	return &DoubleBaseBroker{
		DB:  pool,
		Net: messenger,
	}, nil
}
