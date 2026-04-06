package broker

import (
	"encoding/json"
	"net/http"
)

// HandleGetUsers returns a handler that fetches all users via the broker
func (b *TripleBaseBroker) HandleGetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		data, err := b.Dispatch(r.Context(), "SELECT * FROM users", "rest_api_users")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

// HandleDispatch returns a handler for generic SQL dispatching via the broker
func (b *TripleBaseBroker) HandleDispatch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			SQL     string `json:"sql"`
			Channel string `json:"channel"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		data, err := b.Dispatch(r.Context(), req.SQL, req.Channel)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}
