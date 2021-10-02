package public

import (
	"github.com/gofiber/fiber"
)

type Service interface {
	HostnameAndPort() (string, string)

	SignUp(ctx *fiber.Ctx, req *SignUpRequest) (*SignUpResponse, error)
	Login(ctx *fiber.Ctx, req *LoginRequest) (*LoginResponse, error)
}
