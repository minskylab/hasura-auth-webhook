package server

import (
	"context"

	"github.com/minskylab/hasura-auth-webhook/services/structures"
)

type Service interface {
	// @kok(op): POST /signup
	// @kok(body): body
	// @kok(success): body:res
	PostSignup(ctx context.Context, body structures.PostSignupReq) (res *structures.PostSignupRes, err error)

	// @kok(op): POST /login
	// @kok(body): body
	// @kok(success): body:res
	PostLogin(ctx context.Context, body structures.PostLoginReq) (res *structures.PostLoginRes, err error)

	// @kok(op): POST /webhook
	// @kok(body): body
	// @kok(success): body:res
	PostWebhook(ctx context.Context, body structures.PostWebhookReq) (res *structures.PostWebhookRes, err error)
}
