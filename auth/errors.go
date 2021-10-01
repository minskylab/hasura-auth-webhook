package auth

import (
	"github.com/pkg/errors"
)

var ErrInvalidClaims = errors.New("invalid claims")

var ErrInvalidToken = errors.New("invalid token payload")

var ErrInvalidSignMethod = errors.New("unexpected signing method")

func ErrTokenExpired(e error) error {
	return errors.WithMessage(e, "token is expired")
}
