package server

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/minskylab/hasura-auth-webhook/auth"
	"github.com/minskylab/hasura-auth-webhook/config"
	"github.com/minskylab/hasura-auth-webhook/ent"
	"github.com/minskylab/hasura-auth-webhook/ent/user"
	"github.com/minskylab/hasura-auth-webhook/helpers"
	"github.com/minskylab/hasura-auth-webhook/services"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

const (
	refreshTokenCookieName = "refresh-token"
)

type PublicServer struct {
	client *ent.Client
	auth   *auth.AuthManager

	hostname string
	port     int
}

func NewPublicServer(client *ent.Client, auth *auth.AuthManager, conf *config.Config) *PublicServer {
	return &PublicServer{
		client: client,
		auth:   auth,

		hostname: conf.API.Public.Hostname,
		port:     conf.API.Public.Port,
	}
}

func (p *PublicServer) Hostname() string {
	return p.hostname
}

func (p *PublicServer) Port() int {
	return p.port
}

func (p *PublicServer) Register(ctx *fiber.Ctx) error {
	var req *services.SignUpRequest
	if err := ctx.BodyParser(req); err != nil {
		return errorResponse(ctx.Status(400), errors.Wrap(err, "error on parse body"))
	}

	// TODO: role allowed to register new users
	// check header, get role from token
	// valid roles to this endpoint := []
	// if role not in valid_roles
	// error 403

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
		SetEmail(req.Email).
		SetHashedPassword(hashed).
		Save(ctx.Context())
	if err != nil {
		return errorResponse(ctx.Status(500), errors.New("user could not be created"))
	}

	res := services.SignUpResponse{
		UserID: u.ID.String(),
	}

	return ctx.Status(201).JSON(res)
}

func (p *PublicServer) Login(ctx *fiber.Ctx) error {
	var req services.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return errorResponse(ctx.Status(400), errors.Wrap(err, "error on body parser"))
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
		return errorResponse(ctx.Status(400), errors.WithMessage(err, "wrong credentials"))
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
	res := services.LoginResponse{
		UserID:      u.ID.String(),
		AccessToken: accessToken,
	}

	return ctx.Status(200).JSON(res)
}

func (p *PublicServer) RefreshToken(ctx *fiber.Ctx) error {
	// get cookie
	refreshToken := ctx.Cookies(refreshTokenCookieName)

	if refreshToken == "" {
		return errorResponse(ctx.Status(401), errors.New("header not found"))
	}
	receivedRefreshToken := strings.TrimSpace(refreshToken)

	// validate token and get raw data
	payload, err := p.auth.ValidateRefreshToken(receivedRefreshToken)
	if err != nil {
		return errorResponse(ctx.Status(401), errors.Wrap(err, "invalid refresh token"))
	}

	// find user by its userID
	uid, err := uuid.FromString(payload.UserID)
	if err != nil {
		return errorResponse(ctx.Status(401), errors.Wrap(err, "invalid access token"))
	}

	u, err := p.client.User.Query().Where(user.ID(uid)).First(ctx.Context())
	if err != nil {
		return errorResponse(ctx.Status(401), errors.Wrap(err, "user not found or not exist"))
	}

	// generate access token for user
	accessToken, err := p.auth.DispatchAccessToken(u)
	if err != nil {
		return errorResponse(ctx.Status(500), errors.New("there was an error creating the access token"))
	}

	// // generate refresh token for user
	// refreshToken, err = p.auth.DispatchRefreshToken(u)
	// if err != nil {
	// 	return errorResponse(ctx.Status(500), errors.New("there was an error creating the access token"))
	// }

	// ctx.Cookie(&fiber.Cookie{
	// 	Name:  refreshTokenCookieName,
	// 	Value: refreshToken,
	// })

	// parse json response
	res := services.RefreshTokenResponse{
		UserID:      u.ID.String(),
		AccessToken: accessToken,
	}

	return ctx.Status(200).JSON(res)
}
