package auth

import (
	"time"

	"github.com/minskylab/hasura-auth-webhook/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// RawSecret returns a raw secret, use with caution.
func RawSecret(as []byte, rs []byte) SecretSource {
	return rawSecret{accessSecret: as, refreshSecret: rs}
}

// New creates a new AuthManager instance.
func New(source SecretSource, conf *config.Config) (*AuthManager, error) { // anonymous *AnonymousRole
	as, err := source.GetAccessSecret()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rs, err := source.GetRefreshSecret()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	accessTokenDurationString := conf.Providers.Email.JWT.AccessTokenDuration
	refreshTokenDurationString := conf.Providers.Email.JWT.RefreshTokenDuration

	accessTokenDuration, err := time.ParseDuration(accessTokenDurationString)
	if err != nil {
		logrus.Warn("error at parsing access token duration, using default value")
		accessTokenDuration = defaultAccessTokenDuration
	}

	refreshTokenDuration, err := time.ParseDuration(refreshTokenDurationString)
	if err != nil {
		logrus.Warn("error at parsing refresh token duration, using default value")
		refreshTokenDuration = defaultRefreshTokenDuration
	}

	return &AuthManager{
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
		accessSecret:         as,
		refreshSecret:        rs,
		// anonymous:            anonymous,
	}, nil
}
