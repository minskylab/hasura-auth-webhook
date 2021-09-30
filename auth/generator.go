package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	// "github.com/minskylab/fairpay-sdk/ent"
	"github.com/pkg/errors"
)

const issuer = "fairpay.v1"

func anonymousToken(secret []byte, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "00000000000000000000000000",
		"org": "00000000000000000000000000",
		"iss": issuer,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(duration).Unix(),
		"hasura": map[string]interface{}{
			"x-hasura-allowed-roles": []string{"anonymous"},
			"x-hasura-default-role":  "anonymous",
			"x-hasura-user-id":       "00000000000000000000000000",
			"x-hasura-user-name":     "Anonymous",
			"x-hasura-role-id":       "00000000000000000000000000",
		},
	})

	return token.SignedString(secret)
}

func newToken(forUser *ent.User, role *ent.Role, secret []byte, duration time.Duration, otherRoles ...*ent.Role) (string, error) {
	allowedRoles := []string{role.Kind.String()}

	for _, r := range otherRoles {
		if r.Kind.String() == role.Kind.String() {
			continue
		}

		allowedRoles = append(allowedRoles, r.Kind.String())
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": forUser.ID,
		"iss": issuer,
		"org": forOrganization.ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(duration).Unix(),
		"hasura": map[string]interface{}{
			"x-hasura-allowed-roles": allowedRoles,
			"x-hasura-default-role":  role.Kind,
			"x-hasura-user-id":       forUser.ID,
			"x-hasura-user-name":     forUser.Name,
			"x-hasura-role-id":       role.ID,
		},
	})

	return token.SignedString(secret)
}

func validateToken(token string, secret []byte) (*TokenPayload, error) {
	finalToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})

	mapClaims, ok := finalToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	if err != nil {
		if !mapClaims.VerifyExpiresAt(time.Now().Unix(), true) { // not expired
			return nil, errors.WithMessage(err, "token is expired")
		}
	}

	mapClaims, ok = finalToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	userID, ok := mapClaims["sub"].(string)
	if !ok {
		return nil, errors.New("invalid token payload")
	}

	organizationID, ok := mapClaims["org"].(string)
	if !ok {
		return nil, errors.New("invalid token payload")
	}

	roleID, ok := mapClaims["rol"].(string)
	if !ok {
		return nil, errors.New("invalid token payload")
	}

	return &TokenPayload{
		UserID:         userID,
		OrganizationID: organizationID,
		RoleID:         roleID,
	}, nil
}
