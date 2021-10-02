package auth

import "time"

const issuer = "hasura-app"

var (
	defaultAccessTokenDuration  = 15 * time.Minute
	defaultRefreshTokenDuration = 7 * 24 * time.Hour
)
