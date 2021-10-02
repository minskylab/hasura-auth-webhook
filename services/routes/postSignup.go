package routes

import (
	"encoding/json"
	"net/http"

	"github.com/minskylab/hasura-auth-webhook/helpers"
	"github.com/minskylab/hasura-auth-webhook/server"
	"github.com/minskylab/hasura-auth-webhook/services/structures"
)

func (s service) PostSignup(w http.ResponseWriter, r *http.Request) {
	// input validation body
	var req structures.PostSignupReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		server.ResponseError(w, 400, "Wrong body")
		return
	}

	if ok := helpers.ValidateEmail(req.Email); !ok {
		server.ResponseError(w, 400, "Wrong body")
		return
	}

	if ok := helpers.ValidatePassword(req.Password); !ok {
		server.ResponseError(w, 400, "Wrong body")
		return
	}

	// create user on DB
	hashed, err := helpers.HashPassword(req.Password)
	if err != nil {
		server.ResponseError(w, 500, "User could not be created")
	}

	user, err := s.engine.Client.User.Create().
		SetEmail(req.Email).SetHashedPassword(hashed).
		Save(r.Context())

	if err != nil {
		server.ResponseError(w, 500, "User could not be created")
	}

	// parse response
	res := structures.PostSignupRes{
		UserID: user.ID.String(),
	}

	server.ResponseJSON(w, 201, res)
}
