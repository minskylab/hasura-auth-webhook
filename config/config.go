package config

type JWT struct {
	AccessSecret         string         `yaml:"accessSecret"`
	RefreshSecret        string         `yaml:"refreshSecret"`
	RefreshOptions       RefreshOptions `yaml:"refreshOptions"`
	AccessTokenDuration  string         `yaml:"accessTokenDuration"`
	RefreshTokenDuration string         `yaml:"refreshTokenDuration"`
}

type Refresh struct {
	Name     string `yaml:"name"`
	Domain   string `yaml:"domain"`
	Secure   bool   `yaml:"secure"`
	HttpOnly bool   `yaml:"httpOnly"`
}

type DB struct {
	URL string `yaml:"url"`
}

type JWT struct {
	AccessSecret  string `yaml:"accessSecret"`
	RefreshSecret string `yaml:"refreshSecret"`
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

type Admin struct {
	Users []User `yaml:"users,flow,omitempty"`
}

type Anonymous struct {
	Name string `yaml:"name"`
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

<<<<<<< Updated upstream
type User struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
=======
type EmailProvider struct {
	Enabled  bool          `yaml:"enabled"`
	JWT      JWT           `yaml:"jwt"`
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
>>>>>>> Stashed changes
}

type UserInfo struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

type Mailersend struct {
	Key      string   `yaml:"key"`
	Template string   `yaml:"template"`
	Support  string   `yaml:"support"`
	Name     string   `yaml:"name"`
	User     UserInfo `yaml:"user"`
	Url      string   `yaml:"url"`
}
