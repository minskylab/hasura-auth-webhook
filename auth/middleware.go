package auth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// FairpayContextKey ...
type FairpayContextKey string

const AuthContextKey FairpayContextKey = "auth-context"

// const introspectionSecretBearer = "introspectionSecret"

// AuthenticatorOptions ...
type AuthenticatorOptions struct {
	ProductionMode      bool
	IntrospectionSecret string
	// BearerName string
}

// Authenticator ...
func (auth *FairpayAuth) Authenticator(opts *AuthenticatorOptions) gin.HandlerFunc {
	// prod := opts.ProductionMode
	// introspectionSecret := opts.IntrospectionSecret

	return func(c *gin.Context) {
		// if !prod {
		// 	c.Next()
		// 	return
		// }

		// secret := c.DefaultQuery(introspectionSecretBearer, "")

		// if secret != "" && secret == introspectionSecret {
		// 	c.Next()
		// 	return
		// }

		// defaultRole := c.Request.Header.Get(defaultRoleHeaderName)
		// defaultRoleX := c.Request.Header.Get("X-Hasura-Role")

		// logrus.Info(defaultRole, defaultRoleX)

		// if defaultRole != "manager" {
		// 	c.Status(http.StatusUnauthorized)
		// 	logrus.Errorf("invalid hasura role: '%s'", defaultRole)
		// 	c.Next()
		// 	return
		// }

		authContext, err := contextFromRequest(c.Request)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			logrus.Error(err)
			c.Next()
			return
		}

		ctx := context.WithValue(c.Request.Context(), AuthContextKey, authContext)

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
