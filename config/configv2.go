package config

type JWT2 struct {
	AccessSecret   string  `yaml:"accessSecret"`
	RefreshSecret  string  `yaml:"refreshSecret"`
	RefreshOptions Refresh `yaml:"refreshOptions"`
}

type Auth struct {
	API   API   `yaml:"api"`
	JWT   JWT2  `yaml:"jwt"`
	Admin Admin `yaml:"admin"`
}

type Config2 struct {
	Version string `yaml:"version"`

	Database Database `yaml:"database"`
	Auth     Auth     `yaml:"auth"`
	Roles    []Role   `yaml:"roles,mapstructure"`
}

type Webhooks struct {
	RecoveryPasswordEmail Webhook `yaml:"recoveryPasswordEmail"`
	MagicLinkEmail        Webhook `yaml:"magicLinkEmail"`
}

type Webhook struct {
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers"`
}
