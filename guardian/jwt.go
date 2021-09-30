package guardian

import (
	"net/http"

	"github.com/shaj13/go-guardian/v2/auth"
	"github.com/shaj13/go-guardian/v2/auth/strategies/jwt"
	"github.com/shaj13/libcache"
)

var keeper jwt.SecretsKeeper

func getJWTStrategy(cache *libcache.Cache) (auth.Strategy, error) {
	keeper = jwt.StaticSecret{
		ID:        "secret-id",
		Secret:    []byte("secret"),
		Algorithm: jwt.HS256,
	}

	jwtStrategy := jwt.New(*cache, keeper)
	return jwtStrategy, nil
}

func createToken(r *http.Request) (string, error) {
	u := auth.User(r)
	token, _ := jwt.IssueAccessToken(u, keeper)
	return token, nil
}
