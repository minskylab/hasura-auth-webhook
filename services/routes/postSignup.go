package routes

import (
	"encoding/json"
	"net/http"

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

	// create user on DB
	user, err := s.engine.Client.User.Create().
		SetEmail(req.Email).SetHashedPassword(req.Password).
		Save(r.Context())

	if err != nil {
		server.ResponseError(w, 500, "User could not be created")
	}

	// generate token for user
	// u := auth.User(r)
	// token, _ := jwt.IssueAccessToken(u, keeper)

	// parse response
	res := structures.PostSignupRes{
		UserID: user.ID.String(),
	}

	server.ResponseJSON(w, 201, res)
}
