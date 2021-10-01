package auth

import (
	"time"

	"github.com/minskylab/hasura-auth-webhook/ent"
)

// TokenPayload contains the basic three ids of any fairpay token
type TokenPayload struct {
	UserID string
}

// AuthManager provide methods to dispatch and validate tokens.
type AuthManager struct {
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	accessSecret         []byte
	refreshSecret        []byte
}

// DispatchToken use to dispatch a simple jwt token.
func (auth *AuthManager) DispatchAccessToken(forUser *ent.User, duration ...time.Duration) (string, error) {
	if len(duration) == 0 {
		return newToken(forUser.ID.String(), auth.accessSecret, auth.accessTokenDuration)
	}
	return newToken(forUser.ID.String(), auth.accessSecret, duration[0])
}

func (auth *AuthManager) DispatchRefreshToken(forUser *ent.User, duration ...time.Duration) (string, error) {
	if len(duration) == 0 {
		return newToken(forUser.ID.String(), auth.refreshSecret, auth.refreshTokenDuration)
	}
	return newToken(forUser.ID.String(), auth.refreshSecret, duration[0])

}

// ValidateToken validates an existent jwt token.
func (auth *AuthManager) ValidateAccessToken(token string) (*TokenPayload, error) {
	return validateToken(token, auth.accessSecret)
}

// ValidateToken validates an existent jwt token.
func (auth *AuthManager) ValidateRefreshToken(token string) (*TokenPayload, error) {
	return validateToken(token, auth.refreshSecret)
}
