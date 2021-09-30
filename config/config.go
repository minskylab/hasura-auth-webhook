package config

import (
	"os"
)

type Config struct {
	HOST                   string
	PORT                   string
	DB_HOST                string
	DB_PORT                string
	DB_USER                string
	DB_PASS                string
	DB_DATABASE            string
	JWT_ACCESS_KEY_SECRET  string
	JWT_REFRESH_KEY_SECRET string
}

// func getEnv(key string, fallback string) string {
// 	value := os.Getenv(key)
// 	if len(value) == 0 {
// 		return fallback
// 	}
// 	return value
// }

func NewConfig() *Config {
	keys := map[string]string{
		"HOST":                   "0.0.0.0",
		"PORT":                   "8000",
		"DB_HOST":                "",
		"DB_PORT":                "",
		"DB_USER":                "",
		"DB_PASS":                "",
		"DB_DATABASE":            "",
		"JWT_ACCESS_KEY_SECRET":  "",
		"JWT_REFRESH_KEY_SECRET": "",
	}

	// keys["HOST"] = getEnv("HOST", "")
	if value, ok := os.LookupEnv("HOST"); ok {
		keys["HOST"] = value
	}

	if value, ok := os.LookupEnv("PORT"); ok {
		keys["PORT"] = value
	}

	if value, ok := os.LookupEnv("DB_HOST"); ok {
		keys["DB_HOST"] = value
	}

	if value, ok := os.LookupEnv("DB_PORT"); ok {
		keys["DB_PORT"] = value
	}

	if value, ok := os.LookupEnv("DB_USER"); ok {
		keys["DB_USER"] = value
	}

	if value, ok := os.LookupEnv("DB_PASS"); ok {
		keys["DB_PASS"] = value
	}

	if value, ok := os.LookupEnv("DB_DATABASE"); ok {
		keys["DB_DATABASE"] = value
	}

	if value, ok := os.LookupEnv("JWT_ACCESS_KEY_SECRET"); ok {
		keys["JWT_ACCESS_KEY_SECRET"] = value
	}

	if value, ok := os.LookupEnv("JWT_REFRESH_KEY_SECRET"); ok {
		keys["JWT_REFRESH_KEY_SECRET"] = value
	}

	return &Config{
		HOST:                   keys["HOST"],
		PORT:                   keys["PORT"],
		DB_HOST:                keys["DB_HOST"],
		DB_PORT:                keys["DB_PORT"],
		DB_USER:                keys["DB_USER"],
		DB_PASS:                keys["DB_PASS"],
		DB_DATABASE:            keys["DB_DATABASE"],
		JWT_ACCESS_KEY_SECRET:  keys["JWT_ACCESS_KEY_SECRET"],
		JWT_REFRESH_KEY_SECRET: keys["JWT_REFRESH_KEY_SECRET"],
	}
}
