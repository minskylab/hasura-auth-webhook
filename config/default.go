package config

func NewDefaultConfig(roles ...Role) *Config {
	config := &Config{}

	config.API.Internal.Hostname = "0.0.0.0"
	config.API.Internal.Port = 1111

	config.API.Public.Hostname = "0.0.0.0"
	config.API.Public.Port = 8080

	config.Admin.Users = append(config.Admin.Users, User{
		Email: "admin@example.com", Password: "no-secure-password",
	})

	config.DB.URL = "sqlite:auth.db?loc=auto&_fk=1"

	config.JWT.AccessSecret = "changemeplease"
	config.JWT.RefreshSecret = "changemeplease"

	return config
}