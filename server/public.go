package server

import (
	"net/mail"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/gofiber/fiber/v2"
	mailersend "github.com/mailersend/mailersend-go"
	cache "github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/minskylab/hasura-auth-webhook/auth"
	"github.com/minskylab/hasura-auth-webhook/config"
	"github.com/minskylab/hasura-auth-webhook/ent"
	"github.com/minskylab/hasura-auth-webhook/ent/user"
	"github.com/minskylab/hasura-auth-webhook/helpers"
	"github.com/minskylab/hasura-auth-webhook/services"
)

type RefreshCookie struct {
	Name     string `yaml:"name"`
	Domain   string `yaml:"domain"`
	Secure   bool   `yaml:"secure"`
	HttpOnly bool   `yaml:"httpOnly"`
}

type PublicServer struct {
	Client *ent.Client
	Auth   *auth.AuthManager
	cache  *cache.Cache

	hostname string
	port     int

	Refresh    *RefreshCookie
	Mailersend *mailersend.Mailersend
	Config     *config.Config
}

func NewPublicServer(client *ent.Client, auth *auth.AuthManager, conf *config.Config) services.PublicService {
	return &PublicServer{
		Client: client,
		Auth:   auth,
		// cache: cache.New(24*time.Hour, 12*time.Hour),
		cache: cache.New(5*time.Second, 10*time.Second),

		hostname: conf.API.Public.Hostname,
		port:     conf.API.Public.Port,

		Refresh:    (*RefreshCookie)(conf.Refresh),
		Mailersend: mailersend.NewMailersend(conf.Mailersend.Key),
		Config:     conf,
	}
}

func (p *PublicServer) Hostname() string {
	return p.hostname
}

func (p *PublicServer) Port() int {
	return p.port
}

func (p *PublicServer) Register(ctx *fiber.Ctx) error {
	var req services.SignUpRequest
	if err := ctx.BodyParser(&req); err != nil {
		return errorResponse(ctx.Status(400), errors.Wrap(err, "error on parse body"))
	}

	// TODO: role allowed to register new users
	// check header, get role from token
	authorizationHeader := ctx.Get(authorizationHeaderName)
	withBearerToken := strings.HasPrefix(authorizationHeader, bearerTokenWord)

	if !withBearerToken {
		return errorResponse(ctx.Status(401), errors.New("header not found"))
	}

	receivedAccessToken := strings.TrimSpace(strings.ReplaceAll(authorizationHeader, bearerTokenWord, ""))

	// validate token and get raw data
	payload, err := p.Auth.ValidateAccessToken(receivedAccessToken)
	if err != nil {
		return errorResponse(ctx.Status(401), errors.Wrap(err, "invalid access token"))
	}

	// find user and roles by its userID
	uid, err := uuid.FromString(payload.UserID)
	if err != nil {
		return errorResponse(ctx.Status(401), errors.Wrap(err, "invalid access token"))
	}

	me, err := p.Client.User.Query().Where(user.ID(uid)).First(ctx.Context())
	if err != nil {
		return errorResponse(ctx.Status(401), errors.Wrap(err, "user not found or not exist"))
	}

	myRoles, err := me.QueryRoles().All(ctx.Context())
	if err != nil {
		return errorResponse(ctx.Status(401), errors.Wrap(err, "user hasn't valid roles"))
	}

	// valid roles to this endpoint := []
	// if role not in valid_roles
	// error 403
	if !roleInRoles(myRoles, "admin", "doctor") {
		return errorResponse(ctx.Status(403), errors.New("user hasn't enough permissions"))
	}

	if ok := helpers.ValidateEmail(req.Email); !ok {
		return errorResponse(ctx.Status(400), errors.New("wrong input data"))
	}

	if ok := helpers.ValidatePassword(req.Password); !ok {
		return errorResponse(ctx.Status(400), errors.New("wrong input data"))
	}

	hashed, err := helpers.HashPassword(req.Password)
	if err != nil {
		return errorResponse(ctx.Status(500), errors.Wrap(err, "user could not be created"))
	}

	u, err := p.Client.User.Create().
		SetEmail(req.Email).
		SetHashedPassword(hashed).
		Save(ctx.Context())
	if err != nil {
		return errorResponse(ctx.Status(500), errors.Wrap(err, "user could not be created"))
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
	u, err := p.Client.User.Query().Where(user.Email(req.Email)).Only(ctx.Context())
	if err != nil {
		return errorResponse(ctx.Status(400), errors.WithMessage(err, "wrong credentials"))
	}

	// compare password
	if ok := helpers.CheckPasswordHash(req.Password, u.HashedPassword); !ok {
		return errorResponse(ctx.Status(400), errors.New("wrong credentials"))
	}

	// generate access token for user
	accessToken, err := p.Auth.DispatchAccessToken(u)
	if err != nil {
		return errorResponse(ctx.Status(500), errors.New("there was an error creating the access token"))
	}

	// generate Refresh token for user
	refreshToken, err := p.Auth.DispatchRefreshToken(u)
	if err != nil {
		return errorResponse(ctx.Status(500), errors.New("there was an error creating the access token"))
	}

	cookieOps := fiber.Cookie{
		Name:  defaultRefreshTokenCookieName,
		Value: refreshToken,
	}

	if p.Refresh != nil {
		cookieOps.Name = p.Refresh.Name
		cookieOps.Domain = p.Refresh.Domain
		cookieOps.HTTPOnly = p.Refresh.HttpOnly
		cookieOps.Secure = p.Refresh.Secure
	}

	ctx.Cookie(&cookieOps)

	// parse json response
	res := services.LoginResponse{
		UserID:      u.ID.String(),
		AccessToken: accessToken,
	}

	return ctx.Status(200).JSON(res)
}

func (p *PublicServer) RefreshToken(ctx *fiber.Ctx) error {
	// get cookie
	cookieName := defaultRefreshTokenCookieName

	if p.Refresh != nil {
		cookieName = p.Refresh.Name
	}

	refreshToken := ctx.Cookies(cookieName)

	if refreshToken == "" {
		return errorResponse(ctx.Status(401), errors.New("header not found"))
	}
	receivedRefreshToken := strings.TrimSpace(refreshToken)

	// validate token and get raw data
	payload, err := p.Auth.ValidateRefreshToken(receivedRefreshToken)
	if err != nil {
		return errorResponse(ctx.Status(401), errors.Wrap(err, "invalid Refresh token"))
	}

	// find user by its userID
	uid, err := uuid.FromString(payload.UserID)
	if err != nil {
		return errorResponse(ctx.Status(401), errors.Wrap(err, "invalid access token"))
	}

	u, err := p.Client.User.Query().Where(user.ID(uid)).First(ctx.Context())
	if err != nil {
		return errorResponse(ctx.Status(401), errors.Wrap(err, "user not found or not exist"))
	}

	// generate access token for user
	accessToken, err := p.Auth.DispatchAccessToken(u)
	if err != nil {
		return errorResponse(ctx.Status(500), errors.New("there was an error creating the access token"))
	}

	// // generate Refresh token for user
	// refreshToken, err = p.Auth.DispatchRefreshToken(u)
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

func (p *PublicServer) RecoverPassword(ctx *fiber.Ctx) error {
	email := string(ctx.Request().URI().QueryArgs().Peek("email"))
	name := string(ctx.Request().URI().QueryArgs().Peek("name"))
	if email == "" {
		return errorResponse(ctx.Status(fiber.StatusBadRequest), errors.New("error: missing email"))
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return errorResponse(ctx.Status(fiber.StatusBadRequest), errors.New("error: invalid email"))
	}

	u, err := p.Client.User.Query().Where(user.Email(email)).Only(ctx.Context())
	if err != nil {
		return errorResponse(ctx.Status(fiber.StatusBadRequest), errors.WithMessage(err, "wrong credentials"))
	}
	type claims struct {
		jwt.StandardClaims
		Email string `json:"email"`
	}

	key, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
		Email: email,
	}).SignedString([]byte(uuid.NewV4().String()))
	if err != nil {
		return errorResponse(ctx.Status(fiber.StatusInternalServerError), err)
	}

	{
		err = p.Client.User.UpdateOne(u).SetRecoverPasswordToken(key).Exec(ctx.Context())
		if err != nil {
			return err
		}
	}

	from := mailersend.From{
		Name:  p.Config.Mailersend.User.Name,
		Email: p.Config.Mailersend.User.Email,
	}

	recipients := []mailersend.Recipient{
		{
			Name:  name,
			Email: email,
		},
	}

	personalization := []mailersend.Personalization{
		{
			Email: email,
			Data: map[string]interface{}{
				"name":                  name,
				"recovery_redirect_url": p.Config.Mailersend.Url + key,
				"support_email":         p.Config.Mailersend.Support,
				"account_name":          p.Config.Mailersend.Name,
			},
		},
	}

	message := p.Mailersend.Email.NewMessage()

	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject("Password Recovery")
	message.SetTemplateID(p.Config.Mailersend.Template)
	message.SetPersonalization(personalization)

	_, err = p.Mailersend.Email.Send(ctx.Context(), message)
	if err != nil {
		return errorResponse(ctx.Status(fiber.StatusInternalServerError), err)
	}

	return nil
}

func (p *PublicServer) ChangePassword(ctx *fiber.Ctx) error {
	type claims struct {
		jwt.StandardClaims
		Email string `json:"email"`
	}
	var req services.ChangePasswordRequest
	if err := ctx.BodyParser(&req); err != nil {
		return errorResponse(ctx.Status(fiber.StatusBadRequest), errors.Wrap(err, "error on body parser"))
	}

	email := ""
	var err error = nil

	_, _ = jwt.ParseWithClaims(req.Token, &claims{}, func(token *jwt.Token) (interface{}, error) {
		email = token.Claims.(*claims).Email
		err = token.Claims.Valid()
		return "", nil
	})
	if err != nil {
		return errorResponse(ctx.Status(fiber.StatusBadRequest), err)
	}

	if email == "" {
		return errorResponse(ctx.Status(fiber.StatusBadRequest), errors.New("invalid token"))
	}

	u, err := p.Client.User.Query().Where(user.Email(email)).Only(ctx.Context())
	if err != nil {
		return errorResponse(ctx.Status(fiber.StatusBadRequest), errors.WithMessage(err, "wrong credentials"))
	}

	if u.RecoverPasswordToken != req.Token {
		return errorResponse(ctx.Status(fiber.StatusBadRequest), errors.New("wrong token"))
	}

	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		return errorResponse(ctx.Status(fiber.StatusInternalServerError), err)
	}

	err = p.Client.User.UpdateOne(u).SetHashedPassword(hashedPassword).SetRecoverPasswordToken("").Exec(ctx.Context())
	if err != nil {
		return errorResponse(ctx.Status(fiber.StatusInternalServerError), err)
	}

	return nil
}
