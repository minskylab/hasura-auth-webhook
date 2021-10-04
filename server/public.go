package server

import (
	"github.com/gofiber/fiber"
	"github.com/minskylab/hasura-auth-webhook/auth"
	"github.com/minskylab/hasura-auth-webhook/config"
	"github.com/minskylab/hasura-auth-webhook/ent"
	"github.com/minskylab/hasura-auth-webhook/ent/user"
	"github.com/minskylab/hasura-auth-webhook/helpers"
	"github.com/minskylab/hasura-auth-webhook/services/public"
	"github.com/pkg/errors"
)

type PublicServer struct {
	client *ent.Client
	auth   *auth.AuthManager

	Hostname string
	Port     int
}

func NewPublicServer(client *ent.Client, auth *auth.AuthManager, conf *config.Config) *PublicServer {
	return &PublicServer{
		client: client,
		auth:   auth,

		Hostname: conf.API.Public.Hostname,
		Port:     conf.API.Public.Port,
	}
}

func (p *PublicServer) HostnameAndPort() (string, int) {
	return p.Hostname, p.Port
}

func (p *PublicServer) SignUp(ctx *fiber.Ctx) error {
	var req *public.SignUpRequest
	if err := ctx.BodyParser(req); err != nil {
		return errorResponse(ctx.Status(400), errors.Wrap(err, "error on parse body"))
	}

	if ok := helpers.ValidateEmail(req.Email); !ok {
		return errorResponse(ctx.Status(400), errors.New("wrong input data"))
	}

	if ok := helpers.ValidatePassword(req.Password); !ok {
		return errorResponse(ctx.Status(400), errors.New("wrong input data"))
	}

	hashed, err := helpers.HashPassword(req.Password)
	if err != nil {
		return errorResponse(ctx.Status(500), errors.New("user could not be created"))
	}

	u, err := p.client.User.Create().
		SetEmail(req.Email).SetHashedPassword(hashed).
		Save(ctx.Context())
	if err != nil {
		return errorResponse(ctx.Status(500), errors.New("user could not be created"))
	}

	res := public.SignUpResponse{
		UserID: u.ID.String(),
	}

	return ctx.Status(201).JSON(res)
}

func (p *PublicServer) Login(ctx *fiber.Ctx) error {
	var req *public.LoginRequest
	if err := ctx.BodyParser(req); err != nil {
		return errorResponse(ctx.Status(400), errors.Wrap(err, "error on parse body"))
	}

	if ok := helpers.ValidateEmail(req.Email); !ok {
		return errorResponse(ctx.Status(400), errors.New("wrong input data"))
	}

	if ok := helpers.ValidatePassword(req.Password); !ok {
		return errorResponse(ctx.Status(400), errors.New("wrong input data"))
	}

	// lookup user by email
	u, err := p.client.User.Query().Where(user.Email(req.Email)).Only(ctx.Context())
	if err != nil {
		return errorResponse(ctx.Status(400), errors.New("wrong credentials"))
	}

	// compare password
	if ok := helpers.CheckPasswordHash(req.Password, u.HashedPassword); !ok {
		return errorResponse(ctx.Status(400), errors.New("wrong credentials"))
	}

	// generate access token for user
	accessToken, err := p.auth.DispatchAccessToken(u)
	if err != nil {
		return errorResponse(ctx.Status(500), errors.New("there was an error creating the access token"))
	}

	// generate refresh token for user
	refreshToken, err := p.auth.DispatchRefreshToken(u)
	if err != nil {
		return errorResponse(ctx.Status(500), errors.New("there was an error creating the access token"))
	}

	ctx.Cookie(&fiber.Cookie{
		Name:  "refresh-token",
		Value: refreshToken,
	})

	// parse json response
	res := public.LoginResponse{
		UserID:      u.ID.String(),
		AccessToken: accessToken,
	}

	return ctx.Status(200).JSON(res)
}
