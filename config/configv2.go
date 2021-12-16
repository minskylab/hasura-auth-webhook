package config

type JWT2 struct {
	AccessSecret   string  `yaml:"accessSecret"`
	RefreshSecret  string  `yaml:"refreshSecret"`
	RefreshOptions Refresh `yaml:"refreshOptions"`
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

type Config2 struct {
	Database  Database  `yaml:"database"`
	API       API       `yaml:"api"`
	Admin     Admin     `yaml:"admin"`
	Roles     []Role    `yaml:"roles,mapstructure"`
	Providers Providers `yaml:"providers"`
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

func ConfigV2ToConfigV1(c *Config2) *Config {
	confV1 := new(Config)
	confV1.API = c.API
	confV1.DB = c.Database
	confV1.Admin = c.Admin
	confV1.JWT = JWT{
		AccessSecret:  c.Providers.Email.JWT.AccessSecret,
		RefreshSecret: c.Providers.Email.JWT.RefreshSecret,
	}
	confV1.Refresh = &c.Providers.Email.JWT.RefreshOptions
	confV1.Roles = c.Roles
	confV1.Webhooks = Webhooks{
		Email:     c.Providers.Email.Webhooks,
		MagicLink: c.Providers.MagicLink.Webhooks,
	}

	return confV1
}
