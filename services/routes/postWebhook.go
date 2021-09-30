package routes

import (
	"context"

	"github.com/minskylab/hasura-auth-webhook/services/structures"
)

func (s service) PostWebhook(ctx context.Context, body structures.PostWebhookReq) (*structures.PostWebhookRes, error) {
	return &structures.PostWebhookRes{}, nil
}
