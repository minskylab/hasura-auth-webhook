package helpers

import (
	"context"
	"strings"

	"github.com/minskylab/hasura-auth-webhook/auth"
	"github.com/minskylab/hasura-auth-webhook/ent/user"

	"github.com/minskylab/hasura-auth-webhook/ent"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func ValidateAndGetUserDataFromToken(client *ent.Client, auth *auth.AuthManager, ctx context.Context, authorizationHeader string, bearerTokenWord string) (*ent.User, error) {
	withBearerToken := strings.HasPrefix(authorizationHeader, bearerTokenWord)

	if !withBearerToken {
		return nil, errors.New("header not found")
	}

	receivedAccessToken := strings.TrimSpace(strings.ReplaceAll(authorizationHeader, bearerTokenWord, ""))

	// validate token and get raw data
	payload, err := auth.ValidateAccessToken(receivedAccessToken)
	if err != nil {
		return nil, errors.Wrap(err, "invalid access token")
	}

	// find user and roles by its userID
	uid, err := uuid.FromString(payload.UserID)
	if err != nil {
		return nil, errors.Wrap(err, "invalid access token")
	}

	me, err := client.User.Query().WithRoles().Where(user.ID(uid)).First(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "user not found or not exist")
	}

	return me, nil
}
