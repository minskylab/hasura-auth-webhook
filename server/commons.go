package server

import (
	"context"

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

func searchRolesInParents(ctx context.Context, myRoles []*ent.Role, parentSearchRol *ent.Role) (bool, error) {
	type boolAndError struct {
		result bool
		err    error
	}

	parentRoles, err := parentSearchRol.QueryParents().All(ctx)
	if err != nil {
		return false, err
	}

	for _, p := range parentRoles {
		if roleInRoles(myRoles, p.Name) {
			return true, nil
		}
	}

	responseChannel := make(chan boolAndError)

	result := false
	for _, p := range parentRoles {
		p := p
		go func(c chan boolAndError) {
			if parentSearchRol.ID == p.ID {
				c <- boolAndError{
					result: false,
					err:    err,
				}
				return
			}
			r, err := searchRolesInParents(ctx, myRoles, p)
			res := boolAndError{
				result: r,
				err:    err,
			}
			c <- res
		}(responseChannel)
	}

	for i := 0; i < len(parentRoles); i++ {
		response := <-responseChannel
		if response.err != nil {
			return false, err
		}

		result = result || response.result
	}

	return result, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
