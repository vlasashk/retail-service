package healthz

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

func HealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"status": "ok"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error().Err(err).Send()
	}
}
