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

func roleInRoles(list []*ent.Role, a ...string) bool {
	for _, b := range list {
		if b != nil && contains(a, b.Name) {
			return true
		}
	}
	return false
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}