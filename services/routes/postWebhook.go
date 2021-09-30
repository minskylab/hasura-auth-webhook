package routes

import (
	"encoding/json"
	"net/http"

	"github.com/minskylab/hasura-auth-webhook/server"
	"github.com/minskylab/hasura-auth-webhook/services/structures"
)

func (s service) PostWebhook(w http.ResponseWriter, r *http.Request) {
	// input validation body
	var req structures.PostWebhookReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		server.ResponseError(w, 400, "Wrong body")
		return
	}
	// receivedAccessToken := req.Headers.Bearer

	// validate token and get raw data

	// find user and roles by its userID

	// parse json response
	res := structures.PostWebhookRes{
		HasuraUserId: "",
		HasuraRole:   "",
	}

	server.ResponseJSON(w, 201, res)
}
