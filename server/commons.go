package server

import (
	"github.com/gofiber/fiber/v2"

	"github.com/minskylab/hasura-auth-webhook/ent"
)

const (
	authorizationHeaderName       = "Authorization"
	bearerTokenWord               = "Bearer"
	defaultRefreshTokenCookieName = "Refresh-token"
)

func errorResponse(ctx *fiber.Ctx, err error) error {
	return ctx.JSON(fiber.Map{
		"error": err.Error(),
	})
}

func roleInRoles(a string, list []*ent.Role) bool {
	for _, b := range list {
		if b.Name == a {
			return true
		}
	}
	return false
}
