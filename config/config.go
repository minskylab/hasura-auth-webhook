package config

import (
	"os"
)

type Config struct {
	Host                string
	Port                string
	DBHost              string
	DBPort              string
	DBUser              string
	DBPass              string
	DBDatabase          string
	JwtAccessKeySecret  string
	JwtRefreshKeySecret string
}

func getEnv(key string, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return value
}

func NewConfig() *Config {
	var c Config

	c.Host = getEnv("HOST", "0.0.0.0")
	c.Port = getEnv("PORT", "8000")

	c.DBHost = getEnv("DB_HOST", "")
	c.DBPort = getEnv("DB_PORT", "")
	c.DBUser = getEnv("DB_USER", "")
	c.DBPass = getEnv("DB_PASS", "")
	c.DBDatabase = getEnv("DB_DATABASE", "")

	c.JwtAccessKeySecret = getEnv("JWT_ACCESS_KEY_SECRET", "a-change-me")
	c.JwtRefreshKeySecret = getEnv("JWT_REFRESH_KEY_SECRET", "r-change-me")

	return &c
}
