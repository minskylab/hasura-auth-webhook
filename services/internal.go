package services

import "github.com/gofiber/fiber/v2"

type InternalService interface {
	Hostname() string
	Port() int

	HasuraWebhook(ctx *fiber.Ctx) error
	ListUsers(ctx *fiber.Ctx) error
	// GetUserById(ctx *fiber.Ctx) error
	// CreateUserById(ctx *fiber.Ctx) error
	// UpdateUserById(ctx *fiber.Ctx) error
	// DeleteUserById(ctx *fiber.Ctx) error
}

type HasuraWebhookRequest struct{}

type HasuraWebhookResponse struct {
	HasuraUserId string `json:"X-Hasura-User-Id"`
	HasuraRole   string `json:"X-Hasura-Role"`
}

type Role struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`

	Roles []Role `json:"roles"`
}

type ListUsersResponse []User
