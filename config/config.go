package config

type Config struct {
	API   API    `yaml:"api"`
	DB    DB     `yaml:"db"`
	JWT   JWT    `yaml:"jwt"`
	Roles []Role `yaml:"roles"`
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

type User struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

type Role struct {
	Name  string `yaml:"name"`
	Users []User `yaml:"users,omitempty"`
}
