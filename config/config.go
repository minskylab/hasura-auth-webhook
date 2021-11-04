package config

type Config struct {
	API       API        `yaml:"api"`
	Refresh   *Refresh   `yaml:"refresh"`
	DB        DB         `yaml:"db"`
	JWT       JWT        `yaml:"jwt"`
	Admin     Admin      `yaml:"admin"`
	Anonymous *Anonymous `yaml:"admin"`
	Roles     []Role     `yaml:"roles,mapstructure"`
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
	Name string
}

type Role struct {
	Name  string `yaml:"name"`
	Users []User `yaml:"users,flow,omitempty"`
}

type User struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}
