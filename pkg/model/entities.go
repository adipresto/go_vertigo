package model

import jsoniter "github.com/json-iterator/go"

var jsonInstance = jsoniter.ConfigCompatibleWithStandardLibrary

// User is our declarative entity model
type User struct {
	ID    int64  `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Email string `db:"email" json:"email"`
}


type Payload struct {
	QueryID string `json:"query_id"`
	Weight  int64  `json:"weight"`
	Data    string `json:"data"` // Raw JSON string
}

func WrapPayload(queryID string, weight int64, data []byte) []byte {
	p := Payload{
		QueryID: queryID,
		Weight:  weight,
		Data:    string(data),
	}
	// Using standard JSON here for the wrapper is fine as it's small
	res, _ := jsonInstance.Marshal(p)
	return res
}
