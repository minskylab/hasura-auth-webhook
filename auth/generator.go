package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func newToken(userID string, secret []byte, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": issuer,
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(duration).Unix(),
	})

	return token.SignedString(secret)
}

func validateToken(token string, secret []byte) (*TokenPayload, error) {
	finalToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSignMethod
		}
		return secret, nil
	})

	mapClaims, ok := finalToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	if err != nil {
		if !mapClaims.VerifyExpiresAt(time.Now().Unix(), true) { // not expired
			return nil, ErrTokenExpired(err)
		}
	}

	mapClaims, ok = finalToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	userID, ok := mapClaims["sub"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	iss, ok := mapClaims["iss"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}
	if iss != issuer {
		return nil, ErrInvalidToken
	}

	return &TokenPayload{
		UserID: userID,
	}, nil
}
