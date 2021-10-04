package server

import (
	"strings"

	"github.com/gofiber/fiber"
	"github.com/minskylab/hasura-auth-webhook/auth"
	"github.com/minskylab/hasura-auth-webhook/config"
	"github.com/minskylab/hasura-auth-webhook/ent"
	"github.com/minskylab/hasura-auth-webhook/ent/user"
	"github.com/minskylab/hasura-auth-webhook/services/internal"
	"github.com/minskylab/hasura-auth-webhook/services/structures"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

const (
	authorizationHeaderName = "Authorization"
	bearerTokenWord         = "Bearer"
)

type InternalServer struct {
	client *ent.Client
	auth   *auth.AuthManager

	Hostname string
	Port     int
}

func NewInternalServer(client *ent.Client, auth *auth.AuthManager, conf *config.Config) *InternalServer {
	return &InternalServer{
		client: client,
		auth:   auth,

		Hostname: conf.API.Internal.Hostname,
		Port:     conf.API.Internal.Port,
	}
}

func (s *InternalServer) HostnameAndPort() (string, int) {
	return s.Hostname, s.Port
}

func (s *InternalServer) HasuraWebhook(ctx fiber.Ctx) error {
	// input validation body
	var req *internal.HasuraWebhookRequest
	if err := ctx.BodyParser(req); err != nil {
		return errorResponse(ctx.Status(400), errors.Wrap(err, "error on parse body"))
	}

	authorizationHeader := ctx.Get(authorizationHeaderName)
	// authorizationHeader := r.Header.Get("Authorization")

	if !strings.HasPrefix(authorizationHeader, bearerTokenWord) {
		return errorResponse(ctx.Status(401), errors.New("header not found"))
	}

	receivedAccessToken := strings.TrimSpace(strings.ReplaceAll(authorizationHeader, bearerTokenWord, ""))

	// validate token and get raw data
	payload, err := s.auth.ValidateAccessToken(receivedAccessToken)
	if err != nil {
		return errorResponse(ctx.Status(401), errors.Wrap(err, "invalid access token"))
	}

	// find user and roles by its userID
	uid, err := uuid.FromString(payload.UserID)
	if err != nil {
		return errorResponse(ctx.Status(401), errors.Wrap(err, "invalid access token"))
	}

	u, err := s.client.User.Query().Where(user.ID(uid)).First(ctx.Context())
	if err != nil {
		return errorResponse(ctx.Status(401), errors.Wrap(err, "user not found or not exist"))
	}

	roles, err := u.QueryRoles().All(ctx.Context())
	if err != nil {
		return errorResponse(ctx.Status(401), errors.Wrap(err, "user hasn't valid roles"))
	}

	// parse json response
	res := structures.PostWebhookRes{
		HasuraUserId: "",
		HasuraRole:   roles[0].Name,
	}

	return ctx.Status(201).JSON(res)
}
