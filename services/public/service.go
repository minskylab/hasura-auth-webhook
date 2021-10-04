package public

import (
	"github.com/gofiber/fiber"
)

type Service interface {
	HostnameAndPort() (string, string)

	SignUp(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
}
