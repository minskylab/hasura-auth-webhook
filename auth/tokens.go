package auth

import (
	"time"

	"github.com/minskylab/fairpay-sdk/ent"
	"github.com/pkg/errors"
)

var defaultTokenDuration = 4 * 24 * time.Hour

// TokenPayload contains the basic three ids of any fairpay token
type TokenPayload struct {
	UserID         string
	OrganizationID string
	RoleID         string
}

// FairpayAuth provide methods to dispatch and validate tokens.
type FairpayAuth struct {
	duration time.Duration // func(user *ent.User) time.Duration
	secret   []byte
}

// RawSecret returns a raw secret, use with caution.
func RawSecret(secret []byte) SecretSource {
	return rawSecret{secret: secret}
}

// New create a bew FairpayAuth instance.
func New(source SecretSource) (*FairpayAuth, error) {
	secret, err := source.GetSecret()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &FairpayAuth{
		duration: defaultTokenDuration,
		secret:   secret,
	}, nil
}

// DispatchToken use to dispatch a simple jwt token.
func (auth *FairpayAuth) DispatchToken(forUser *ent.User, inOrganization *ent.Organization, role *ent.Role, extraRoles ...*ent.Role) (string, error) {
	return newToken(forUser, inOrganization, role, auth.secret, auth.duration, extraRoles...)
}

// DispatchAnonymousToken an anonymous jwt token.
func (auth *FairpayAuth) DispatchAnonymousToken(duration ...time.Duration) (string, error) {
	if len(duration) == 0 {
		return anonymousToken(auth.secret, auth.duration)
	}

	return anonymousToken(auth.secret, duration[0])
}

// ValidateToken validates an existent jwt token.
func (auth *FairpayAuth) ValidateToken(token string) (*TokenPayload, error) {
	return validateToken(token, auth.secret)
}
