package server

import "github.com/gofiber/fiber"

func errorResponse(ctx *fiber.Ctx, err error) error {
	return ctx.JSON(fiber.Map{
		"error": err.Error(),
	})
}
