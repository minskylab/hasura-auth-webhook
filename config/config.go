package config

// type Config struct {
// 	Host                string
// 	Port                string
// 	DBHost              string
// 	DBPort              string
// 	DBUser              string
// 	DBPass              string
// 	DBDatabase          string
// 	JwtAccessKeySecret  string
// 	JwtRefreshKeySecret string
// }

// func getEnv(key string, fallback string) string {
// 	value, ok := os.LookupEnv(key)
// 	if !ok {
// 		return fallback
// 	}
// 	return value
// }

// func NewConfig() *Config {
// 	var c Config

// 	c.Host = getEnv("HOST", "0.0.0.0")
// 	c.Port = getEnv("PORT", "8000")

// 	c.DBHost = getEnv("DB_HOST", "")
// 	c.DBPort = getEnv("DB_PORT", "")
// 	c.DBUser = getEnv("DB_USER", "")
// 	c.DBPass = getEnv("DB_PASS", "")
// 	c.DBDatabase = getEnv("DB_DATABASE", "")

// 	c.JwtAccessKeySecret = getEnv("JWT_ACCESS_KEY_SECRET", "a-change-me")
// 	c.JwtRefreshKeySecret = getEnv("JWT_REFRESH_KEY_SECRET", "r-change-me")

// 	return &c
// }

type Config struct {
	API   API    `yaml:"api"`
	Roles []Role `yaml:"roles"`
}

type Public struct {
	Hostname string `yaml:"hostname"`
	Port     int    `yaml:"port"`
}

type Internal struct {
	Hostname string `yaml:"hostname"`
	Port     int    `yaml:"port"`
}

type API struct {
	Public   Public   `yaml:"public"`
	Internal Internal `yaml:"internal"`
}

type User struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

type Role struct {
	Name  string `yaml:"name"`
	Users []User `yaml:"users,omitempty"`
}

func NewDefaultConfig(roles ...Role) *Config {
	config := &Config{}

	config.API.Internal.Hostname = "0.0.0.0"
	config.API.Internal.Port = 8080

	config.API.Public.Hostname = "0.0.0.0"
	config.API.Public.Port = 1111

	config.Roles = append(config.Roles, Role{
		Name: "admin",
		Users: []User{
			{Email: "admin@example.com", Password: "no-secure-password"},
		},
	})

	return config
}
