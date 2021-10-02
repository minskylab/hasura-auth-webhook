package server

import (
	"github.com/gofiber/fiber"
	"github.com/minskylab/hasura-auth-webhook/auth"
	"github.com/minskylab/hasura-auth-webhook/ent"
	"github.com/minskylab/hasura-auth-webhook/services/internal"
)

type InternalServer struct {
	client *ent.Client
	auth   *auth.AuthManager

	Hostname string
	Port     string
}

func (s *InternalServer) HostnameAndPort() (string, string) {
	return s.Hostname, s.Port
}

func (s *InternalServer) HasuraWebhook(ctx fiber.Ctx, req *internal.HasuraWebhookRequest) (*internal.HasuraWebhookResponse, error) {
	return nil, nil
}
