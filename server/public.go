package server

import (
	"net/mail"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/minskylab/hasura-auth-webhook/ent/role"
	"github.com/sirupsen/logrus"

	"github.com/golang-jwt/jwt"

	"github.com/gofiber/fiber/v2"
	// mailersend "github.com/mailersend/mailersend-go"
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
	Client     *ent.Client
	Auth       *auth.AuthManager
	Refresh    *RefreshCookie
	Config     *config.Config
	HTTPCLient *resty.Client
	cache      *cache.Cache

	hostname string
	port     int
}

func NewPublicServer(client *ent.Client, auth *auth.AuthManager, conf *config.Config) services.PublicService {
	return &PublicServer{
		Client:     client,
		Auth:       auth,
		HTTPCLient: resty.New(),
		Refresh:    (*RefreshCookie)(conf.Refresh),
		Config:     conf,

		cache:    cache.New(5*time.Second, 10*time.Second),
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

type NewRegisterEvent struct {
	EventAt time.Time `json:"event_at"`
	UserID  string    `json:"user_id"`
}

func (p *PublicServer) Register(ctx *fiber.Ctx) error {
	var req services.SignUpRequest
	if err := ctx.BodyParser(&req); err != nil {
		return errorResponse(ctx.Status(400), errors.Wrap(err, "error on parse body"))
	}

	rol, err := p.Client.Role.Query().Where(role.Name(req.Role)).Only(ctx.Context())
	if err != nil {
		return errorResponse(ctx.Status(401), err)
	}

	if !rol.Public {
		authorizationHeader := ctx.Get(authorizationHeaderName)

		creator, err := helpers.ValidateAndGetUserDataFromToken(p.Client, p.Auth, ctx.Context(), authorizationHeader, bearerTokenWord)
		if err != nil {
			return errorResponse(ctx.Status(401), err)
		}

		creatorRoles := creator.Edges.Roles

		result, err := searchRolesInParents(ctx.Context(), creatorRoles, rol)
		if err != nil {
			return errorResponse(ctx.Status(401), err)
		}

		if !result {
			return errorResponse(ctx.Status(403), errors.New("user hasn't enough permissions"))
		}
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
		AddRoles(rol).
		Save(ctx.Context())
	if err != nil {
		return errorResponse(ctx.Status(500), errors.Wrap(err, "user could not be created"))
	}

	res := services.SignUpResponse{
		UserID: u.ID.String(),
	}

	// TODO: !Non blocking webhook invocation
	if p.Config.Webhooks.Email.RegisterEvent.URL != "" {
		res, err := p.HTTPCLient.R().
			SetHeaders(p.Config.Webhooks.Email.RegisterEvent.Headers).
			SetBody(NewRegisterEvent{
				EventAt: time.Now(),
				UserID:  u.ID.String(),
			}).
			Post(p.Config.Webhooks.Email.RegisterEvent.URL)
		if err != nil {
			// return errorResponse(ctx.Status(fiber.StatusInternalServerError), err)
			logrus.Warn(err)
		}

		logrus.Info(res.String())
	}

	return ctx.Status(201).JSON(res)
}

type LoginMagicLinkEvent struct {
	EventAt time.Time `json:"event_at"`
	Email   string    `json:"email"`

	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (p *PublicServer) generateTokens(u *ent.User) (string, string, error) {
	// generate access token for user
	accessToken, err := p.Auth.DispatchAccessToken(u)
	if err != nil {
		return "", "", errors.New("there was an error creating the access token")
	}

	// generate Refresh token for user
	refreshToken, err := p.Auth.DispatchRefreshToken(u)
	if err != nil {
		return "", "", errors.New("there was an error creating the access token")
	}

	return accessToken, refreshToken, nil
}

func (p *PublicServer) magicLinkFlow(ctx *fiber.Ctx, u *ent.User) error {
	accessToken, refreshToken, err := p.generateTokens(u)
	if err != nil {
		return errorResponse(ctx.Status(500), err)
	}

	eventAt := time.Now()
	res, err := p.HTTPCLient.R().
		SetHeaders(p.Config.Webhooks.MagicLink.LoginEvent.Headers).
		SetBody(LoginMagicLinkEvent{
			EventAt:      eventAt,
			Email:        u.Email,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}).
		Post(p.Config.Webhooks.MagicLink.LoginEvent.URL)
	if err != nil {
		return errorResponse(ctx.Status(fiber.StatusInternalServerError), err)
	}

	logrus.Info(res.String())

	return ctx.Status(fiber.StatusOK).JSON(map[string]string{
		"event_at": eventAt.Format(time.RFC3339),
		"result":   "ok",
	})
}

func (p *PublicServer) Login(ctx *fiber.Ctx) error {
	var req services.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return errorResponse(ctx.Status(400), errors.Wrap(err, "error on body parser"))
	}

	if ok := helpers.ValidateEmail(req.Email); !ok {
		return errorResponse(ctx.Status(400), errors.New("wrong input data"))
	}

	// lookup user by email
	u, err := p.Client.User.Query().Where(user.Email(req.Email)).Only(ctx.Context())
	if err != nil {
		return errorResponse(ctx.Status(400), errors.WithMessage(err, "wrong credentials"))
	}

	magicLinkWebhookExists := p.Config.Webhooks.MagicLink.LoginEvent.URL != ""
	passwordNotExists := req.Password == nil
	passwordIsBlank := req.Password != nil && *req.Password == ""

	if magicLinkWebhookExists && passwordNotExists {
		return p.magicLinkFlow(ctx, u)
	}

	if magicLinkWebhookExists && passwordIsBlank {
		return p.magicLinkFlow(ctx, u)
	}

	if passwordNotExists {
		return errorResponse(ctx.Status(400), errors.New("wrong input data"))
	}

	if ok := helpers.ValidatePassword(*req.Password); !ok {
		return errorResponse(ctx.Status(400), errors.New("wrong input data"))
	}

	// compare password
	if ok := helpers.CheckPasswordHash(*req.Password, u.HashedPassword); !ok {
		return errorResponse(ctx.Status(400), errors.New("wrong credentials"))
	}

	accessToken, refreshToken, err := p.generateTokens(u)
	if err != nil {
		return errorResponse(ctx.Status(500), err)
	}

	cookieOps := fiber.Cookie{
		Name:  defaultRefreshTokenCookieName,
		Value: refreshToken,
	}

	if p.Refresh != nil {
		if p.Refresh.Name != "" {
			cookieOps.Name = p.Refresh.Name
		}

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

	// parse json response
	res := services.RefreshTokenResponse{
		UserID:      u.ID.String(),
		AccessToken: accessToken,
	}

	return ctx.Status(200).JSON(res)
}

type RecoverPasswordEvent struct {
	EventAt time.Time `json:"event_at"`
	UserID  string    `json:"user_id"`
	Token   string    `json:"token"`
}

func (p *PublicServer) RecoverPassword(ctx *fiber.Ctx) error {
	email := string(ctx.Request().URI().QueryArgs().Peek("email"))
	// name := string(ctx.Request().URI().QueryArgs().Peek("name"))
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

	err = p.Client.User.UpdateOne(u).SetRecoverPasswordToken(key).Exec(ctx.Context())
	if err != nil {
		return errorResponse(ctx.Status(fiber.StatusInternalServerError), err)
	}

	// TODO: !Non blocking webhook invocation
	if p.Config.Webhooks.Email.RecoveryPasswordEvent.URL != "" {
		req, err := p.HTTPCLient.R().
			SetHeaders(p.Config.Webhooks.Email.RecoveryPasswordEvent.Headers).
			SetBody(RecoverPasswordEvent{
				EventAt: time.Now(),
				UserID:  u.ID.String(),
				Token:   key,
			}).
			Post(p.Config.Webhooks.Email.RecoveryPasswordEvent.URL)
		if err != nil {
			// return errorResponse(ctx.Status(fiber.StatusInternalServerError), err)
			logrus.Error(err)
		}

		logrus.Info(req.String())
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
