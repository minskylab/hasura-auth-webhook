package guardian

import (
	"github.com/shaj13/go-guardian/v2/auth"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
	_ "github.com/shaj13/libcache/fifo"
)

func SetupAuth() (*union.Union, error) {

	cache, _ := setupCache()

	var strats []auth.Strategy

	jwtStrategy, _ := getJWTStrategy(cache)
	strats = append(strats, jwtStrategy)

	basicStrategy, _ := getBasicStrategy(cache)
	strats = append(strats, basicStrategy)

	strategy := union.New(strats...)

	return &strategy, nil
}
