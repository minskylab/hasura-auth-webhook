package internal

import "github.com/gofiber/fiber"

type Service interface {
	HostnameAndPort() (string, string)

	HasuraWebhook(ctx fiber.Ctx, req *HasuraWebhookRequest) (*HasuraWebhookResponse, error)
}
