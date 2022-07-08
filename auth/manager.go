package auth

import (
	"time"

	"github.com/minskylab/hasura-auth-webhook/ent"
)

// TokenPayload contains the user id of the token
type TokenPayload struct {
	UserID string
	RoleID *string
}

// AuthManager provide methods to dispatch and validate tokens.
type AuthManager struct {
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	accessSecret         []byte
	refreshSecret        []byte
	// anonymous            *AnonymousRole
}

type AnonymousRole struct {
	Name string
}

// DispatchToken use to dispatch a simple jwt token.
func (auth *AuthManager) DispatchAccessToken(forUser *ent.User) (string, error) {
	return newToken(forUser.ID.String(), auth.accessSecret, auth.accessTokenDuration)
}

func (auth *AuthManager) DispatchRefreshToken(forUser *ent.User) (string, error) {
	return newToken(forUser.ID.String(), auth.refreshSecret, auth.refreshTokenDuration)
}

// ValidateToken validates an existent jwt token.
func (auth *AuthManager) ValidateAccessToken(token string) (*TokenPayload, error) {
	return validateToken(token, auth.accessSecret)
}

// ValidateToken validates an existent jwt token.
func (auth *AuthManager) ValidateRefreshToken(token string) (*TokenPayload, error) {
	return validateToken(token, auth.refreshSecret)
}

// func (auth *AuthManager) IsAnonymousAllowed() bool {
// return auth.anonymous != nil
// }

// func (auth *AuthManager) GetAnonymous() *AnonymousRole {
// 	return auth.anonymous
// }
