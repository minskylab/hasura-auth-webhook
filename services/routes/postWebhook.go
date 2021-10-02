package routes

import (
	"net/http"
	"strings"

	"github.com/minskylab/hasura-auth-webhook/ent/user"
	"github.com/minskylab/hasura-auth-webhook/server"
	"github.com/minskylab/hasura-auth-webhook/services/structures"
	uuid "github.com/satori/go.uuid"
)

func (s service) PostWebhook(w http.ResponseWriter, r *http.Request) {
	// input validation body
	authorizationHeader := r.Header.Get("Authorization")
	arr := strings.Split(authorizationHeader, "Bearer")
	if len(arr) != 2 {
		server.ResponseError(w, 401, "Header not found")
		return
	}
	receivedAccessToken := arr[1]

	// validate token and get raw data
	payload, err := s.engine.Auth.ValidateAccessToken(receivedAccessToken)
	if err != nil {
		server.ResponseError(w, 401, "")
		return
	}

	// find user and roles by its userID
	uid, err := uuid.FromString(payload.UserID)
	if err != nil {
		server.ResponseError(w, 401, "")
		return
	}

	u, err := s.engine.Client.User.Query().Where(user.ID(uid)).First(r.Context())
	if err != nil {
		server.ResponseError(w, 401, "")
		return
	}

	roles, err := u.QueryRoles().All(r.Context())
	if err != nil {
		server.ResponseError(w, 401, "")
		return
	}

	// parse json response
	res := structures.PostWebhookRes{
		HasuraUserId: "",
		HasuraRole:   roles[0].Name,
	}

	server.ResponseJSON(w, 200, res)
}
