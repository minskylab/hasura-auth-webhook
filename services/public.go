package services

import (
	"github.com/gofiber/fiber/v2"
)

type PublicService interface {
	Hostname() string
	Port() int

	Register(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	RefreshToken(ctx *fiber.Ctx) error

	RecoverPassword(ctx *fiber.Ctx) error
	ChangePassword(ctx *fiber.Ctx) error
}

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type SignUpResponse struct {
	UserID string `json:"userID"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`

	Role *string `json:"role"`
}

type LoginResponse struct {
	UserID      string `json:"userID"`
	AccessToken string `json:"accessToken"`
}

type RefreshTokenResponse struct {
	UserID      string `json:"userID"`
	AccessToken string `json:"accessToken"`
}

type ChangePasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}
