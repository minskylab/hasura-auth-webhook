package routes

import (
	"context"

	"github.com/minskylab/hasura-auth-webhook/services/structures"
)

func (s service) PostLogin(ctx context.Context, body structures.PostLoginReq) (*structures.PostLoginRes, error) {

	return &structures.PostLoginRes{}, nil
}
