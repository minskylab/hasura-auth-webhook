package routes

import (
	"net/http"
	"time"

	"github.com/minskylab/hasura-auth-webhook/server"
)

func (s service) GetTime(w http.ResponseWriter, _ *http.Request) {
	server.ResponseJSON(w, 200, map[string]time.Time{
		"time": time.Now(),
	})
}
