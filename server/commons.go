package server

import "github.com/gofiber/fiber/v2"

func errorResponse(ctx *fiber.Ctx, err error) error {
	return ctx.JSON(fiber.Map{
		"error": err.Error(),
	})
}
