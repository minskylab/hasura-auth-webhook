package routes

import (
	"encoding/json"
	"net/http"

	"github.com/minskylab/hasura-auth-webhook/ent/user"
	"github.com/minskylab/hasura-auth-webhook/helpers"
	"github.com/minskylab/hasura-auth-webhook/server"
	"github.com/minskylab/hasura-auth-webhook/services/structures"
)

func (s service) PostLogin(w http.ResponseWriter, r *http.Request) {
	// input validation body
	var req structures.PostLoginReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		server.ResponseError(w, 400, "Wrong body")
		return
	}

	// lookup user by email
	u, err := s.engine.Client.User.Query().Where(user.Email(req.Email)).Only(r.Context())
	if err != nil {
		server.ResponseError(w, 400, "Wrong credentials")
		return
	}

	// compare password
	ok := helpers.CheckPasswordHash(req.Password, u.HashedPassword)
	if !ok {
		server.ResponseError(w, 400, "Wrong credentials")
		return
	}

	// generate access token for user
	accessToken, err := s.engine.Auth.DispatchAccessToken(u)
	if err != nil {
		server.ResponseError(w, 500, "There was an error creating the access token")
		return
	}

	// generate refresh token for user
	refreshToken, err := s.engine.Auth.DispatchRefreshToken(u)
	if err != nil {
		server.ResponseError(w, 500, "There was an error creating the access token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "refresh-token",
		Value: refreshToken,
	})

	// parse json response
	res := structures.PostLoginRes{
		UserID:      u.ID.String(),
		AccessToken: accessToken,
	}

	server.ResponseJSON(w, 201, res)
}
