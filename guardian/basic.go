package guardian

import (
	"context"
	"fmt"
	"net/http"

	"github.com/shaj13/go-guardian/v2/auth"
	"github.com/shaj13/go-guardian/v2/auth/strategies/basic"
	"github.com/shaj13/libcache"
)

func validateUser(ctx context.Context, r *http.Request, userName, password string) (auth.Info, error) {
	// here connect to db or any other service to fetch user and validate it.
	if userName == "admin" && password == "admin" {
		return auth.NewDefaultUser("admin", "1", nil, nil), nil
	}

	return nil, fmt.Errorf("invalid credentials")
}

func getBasicStrategy(cache *libcache.Cache) (auth.Strategy, error) {
	basicStrategy := basic.NewCached(validateUser, *cache)
	return basicStrategy, nil
}
