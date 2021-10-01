package auth

import (
	"github.com/pkg/errors"
)

// RawSecret returns a raw secret, use with caution.
func RawSecret(as []byte, rs []byte) SecretSource {
	return rawSecret{accessSecret: as, refreshSecret: rs}
}

// New creates a new AuthManager instance.
func New(source SecretSource) (*AuthManager, error) {
	as, err := source.GetAccessSecret()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rs, err := source.GetRefreshSecret()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &AuthManager{
		accessTokenDuration:  defaultAccessTokenDuration,
		refreshTokenDuration: defaultRefreshTokenDuration,
		accessSecret:         as,
		refreshSecret:        rs,
	}, nil
}
