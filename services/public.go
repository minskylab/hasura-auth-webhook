package services

import (
	"github.com/gofiber/fiber/v2"
)

type PublicService interface {
	Hostname() string
	Port() int

	Register(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
}

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	UserID string `json:"userID"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID      string `json:"userID"`
	AccessToken string `json:"accessToken"`
}
