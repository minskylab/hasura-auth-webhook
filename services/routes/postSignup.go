package routes

import (
	"context"

	"github.com/minskylab/hasura-auth-webhook/services/structures"
)

func (s service) PostSignup(ctx context.Context, body structures.PostSignupReq) (*structures.PostSignupRes, error) {
	// validar input
	// email := body.Email
	// password := body.Password

	// crear en db
	// user, err := s.engine.Client.User.Create().SetEmail(email).SetPassword(password).Save(ctx)
	// if err != nil {
	// 	return nil, errors.New("not access to this resource")
	// }

	// generar token

	// devolver
	return &structures.PostSignupRes{
		AccessToken:  "",
		UserID:       "",
		RefreshToken: "",
	}, nil
}
