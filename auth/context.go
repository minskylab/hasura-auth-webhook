package auth

import (
	"context"
	"net/http"
)

// AuthContext saves the authorization context provided by the hasura request headers.
type AuthContext struct {
	OrganizationID       string
	OrganizationUsername string
	UserID               string
	RoleID               string
	RoleName             string
}

func contextFromRequest(r *http.Request) (*AuthContext, error) {
	orgID := r.Header.Get(orgIDHeaderName)
	orgUsername := r.Header.Get(orgUsernameHeaderName)
	userID := r.Header.Get(userIDHeaderName)
	roleID := r.Header.Get(roleIDHeaderName)
	defaultRole := r.Header.Get(defaultRoleHeaderName)

	// if orgID == "" || orgUsername == "" {
	// 	return nil, errors.New("organization id or username not found in request context")
	// }

	// if userID == "" {
	// 	return nil, errors.New("user id nor found in request context")
	// }

	// if roleID == "" {
	// 	return nil, errors.New("role id nor found in request context")
	// }

	return &AuthContext{
		OrganizationID:       orgID,
		OrganizationUsername: orgUsername,
		UserID:               userID,
		RoleID:               roleID,
		RoleName:             defaultRole,
	}, nil
}

// ExtractAuthContext extract the auth context structure from context.
func ExtractAuthContext(ctx context.Context) *AuthContext {
	c, isOk := ctx.Value(AuthContextKey).(*AuthContext)
	if !isOk {
		return &AuthContext{}
	}

	return c
}

// ExtractOrganizationID extracts organization id from context.
func ExtractOrganizationID(ctx context.Context) string {
	return ExtractAuthContext(ctx).OrganizationID
}

// ExtractOrganizationUsername ...
func ExtractOrganizationUsername(ctx context.Context) string {
	return ExtractAuthContext(ctx).OrganizationUsername
}
