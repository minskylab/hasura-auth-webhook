package config

type JWT2 struct {
	AccessSecret         string         `yaml:"accessSecret"`
	RefreshSecret        string         `yaml:"refreshSecret"`
	RefreshOptions       RefreshOptions `yaml:"refreshOptions"`
	AccessTokenDuration  string         `yaml:"accessTokenDuration"`
	RefreshTokenDuration string         `yaml:"refreshTokenDuration"`
}

type RefreshOptions struct {
	Name     string `yaml:"name"`
	Domain   string `yaml:"domain"`
	Secure   bool   `yaml:"secure"`
	HttpOnly bool   `yaml:"httpOnly"`
}

type Database struct {
	URL string `yaml:"url"`
}

type User struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

type API struct {
	Public   Public   `yaml:"public"`
	Internal Internal `yaml:"internal"`
}

type Public struct {
	Hostname string `yaml:"hostname"`
	Port     int    `yaml:"port"`
}

type Internal struct {
	Hostname string `yaml:"hostname"`
	Port     int    `yaml:"port"`
}

type Admin struct {
	Users []User `yaml:"users,flow,omitempty"`
}

type Role struct {
	Name     string   `yaml:"name"`
	Parent   *string  `yaml:"parent,omitempty"`
	Parents  []string `yaml:"parents,omitempty"`
	Child    *string  `yaml:"child,omitempty"`
	Children []string `yaml:"children,omitempty"`
	Users    []User   `yaml:"users,flow,omitempty"`
	Public   bool     `yaml:"public,omitempty"`
}

type EmailProvider struct {
	Enabled  bool          `yaml:"enabled"`
	JWT      JWT2          `yaml:"jwt"`
	Webhooks EmailWebhooks `yaml:"webhooks"`
}

type MagicLinkProvider struct {
	Enabled  bool              `yaml:"enabled"`
	Webhooks MagicLinkWebhooks `yaml:"webhooks"`
}

type Providers struct {
	Email     EmailProvider     `yaml:"email"`
	MagicLink MagicLinkProvider `yaml:"magicLink"`
}

type Config struct {
	Database  Database  `yaml:"database"`
	API       API       `yaml:"api"`
	Admin     Admin     `yaml:"admin"`
	Roles     []Role    `yaml:"roles,mapstructure"`
	Providers Providers `yaml:"providers"`
	Webhooks  Webhooks  `yaml:"webhooks"`
}

type Webhooks struct {
	Email     EmailWebhooks     `yaml:"email"`
	MagicLink MagicLinkWebhooks `yaml:"magicLink"`
}

type EmailWebhooks struct {
	RecoveryPasswordEvent Webhook `yaml:"recoveryPasswordEvent"`
	RegisterEvent         Webhook `yaml:"registerEvent"`
}

type MagicLinkWebhooks struct {
	LoginEvent Webhook `yaml:"loginEvent"`
}

type Webhook struct {
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers"`
}

// func ConfigV2ToConfigV1(c *Config2) *Config {
// 	confV1 := new(Config)
// 	confV1.API = c.API
// 	confV1.DB = c.Database
// 	confV1.Admin = c.Admin
// 	confV1.JWT = JWT{
// 		AccessSecret:  c.Providers.Email.JWT.AccessSecret,
// 		RefreshSecret: c.Providers.Email.JWT.RefreshSecret,
// 	}
// 	confV1.Refresh = &c.Providers.Email.JWT.RefreshOptions
// 	confV1.Roles = c.Roles
// 	confV1.Webhooks = Webhooks{
// 		Email:     c.Providers.Email.Webhooks,
// 		MagicLink: c.Providers.MagicLink.Webhooks,
// 	}

// 	return confV1
// }
