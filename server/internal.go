package server

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/minskylab/hasura-auth-webhook/auth"
	"github.com/minskylab/hasura-auth-webhook/config"
	"github.com/minskylab/hasura-auth-webhook/ent"
	"github.com/minskylab/hasura-auth-webhook/ent/user"
	"github.com/minskylab/hasura-auth-webhook/services"
)

type InternalServer struct {
	client *ent.Client
	auth   *auth.AuthManager

	hostname string
	port     int
}

func NewInternalServer(client *ent.Client, auth *auth.AuthManager, conf *config.Config) services.InternalService {
	return &InternalServer{
		client: client,
		auth:   auth,

		hostname: conf.API.Internal.Hostname,
		port:     conf.API.Internal.Port,
	}
}

func (s *InternalServer) Hostname() string {
	return s.hostname
}

func (s *InternalServer) Port() int {
	return s.port
}

func (s *InternalServer) HasuraWebhook(ctx *fiber.Ctx) error {
	authorizationHeader := ctx.Get(authorizationHeaderName)

	anonymous := s.auth.GetAnonymous()
	withAnonymousAllowed := anonymous != nil
	withBearerToken := strings.HasPrefix(authorizationHeader, bearerTokenWord)

	if !withBearerToken && !withAnonymousAllowed {
		return errorResponse(ctx.Status(401), errors.New("header not found"))
	}

	if !withBearerToken && withAnonymousAllowed {
		return ctx.Status(200).JSON(services.HasuraWebhookResponse{
			HasuraUserId: "",
			HasuraRole:   anonymous.Name,
		})
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
	res := services.HasuraWebhookResponse{
		HasuraUserId: u.ID.String(),
		HasuraRole:   roles[0].Name,
	}

	return ctx.Status(200).JSON(res)
}

func (s *InternalServer) ListUsers(ctx *fiber.Ctx) error {
	users, err := s.client.User.Query().WithRoles().All(ctx.Context())
	if err != nil {
		return errorResponse(ctx.Status(401), errors.Wrap(err, "users could not be found"))
	}

	// TODO: parse parse json response
	// res := []services.User{
	// 	services.User {
	// 	}
	// }

	return ctx.Status(200).JSON(users)
}
