package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/qosimmax/storage-api/user"
)

// AddServer is a handler that adds a new server
//
//	POST /api/v1/server
//	Responds: 200, 400, 500
//	Body:
//		type ServerData struct {
//			ID      string `json:"id"`
//			Name    string `json:"name"`
//			Address string `json:"address"`
//		}
//
// The handler should accept an interface(s), and should contain only high level
// business logic.
func AddServer(
	db user.ServerAdder,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var payload user.ServerData
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			handleError(
				w,
				fmt.Errorf("error decoding body: %w", err),
				http.StatusBadRequest,
				true,
			)
			return
		}

		serverID, err := db.AddServer(ctx, payload)
		if err != nil {
			handleError(w, err, http.StatusInternalServerError, true)
			return
		}

		response, _ := json.Marshal(map[string]string{"server_id": serverID})
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(response)

	}
}
