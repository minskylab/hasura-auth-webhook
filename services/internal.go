package services

import "github.com/gofiber/fiber/v2"

type InternalService interface {
	Hostname() string
	Port() int

	HasuraWebhook(ctx *fiber.Ctx) error
}

type HasuraWebhookRequest struct{}

type HasuraWebhookResponse struct {
	HasuraUserId string `json:"X-Hasura-User-Id"`
	HasuraRole   string `json:"X-Hasura-Role"`
}
