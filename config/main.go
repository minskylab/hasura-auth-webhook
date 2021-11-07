package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func NewConfig(path string) (*Config, error) {
	config := NewDefaultConfig()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file %s", err)
	}

	for _, key := range viper.AllKeys() {
		rawValue := viper.Get(key)
		switch value := rawValue.(type) {
		case int:
			// fmt.Println("int: " + key)
			viper.Set(key, value)
		case string:
			// fmt.Println("string: " + key)
			viper.Set(key, os.ExpandEnv(value))
		case []Role:
			// fmt.Println("[]role: " + key)
			for _, v := range value {
				v.Name = os.ExpandEnv(v.Name)
				for _, u := range v.Users {
					u.Email = os.ExpandEnv(u.Email)
					u.Password = os.ExpandEnv(u.Password)
				}
			}
			viper.Set(key, value)
		case []User:
			// fmt.Println("[]user: " + key)
			for _, u := range value {
				u.Email = os.ExpandEnv(u.Email)
				u.Password = os.ExpandEnv(u.Password)
			}
			viper.Set(key, value)
		case []interface{}:
			// fmt.Println("[]default: " + key)
			// fmt.Println(value)

			// for _, v := range value {
			// 	r, ok := v.(Role)
			// 	helpers.Log(r)
			// 	if ok {
			// 		fmt.Printf("\nRole type - %s", r.Name)
			// 		r.Name = os.ExpandEnv(r.Name)
			// 		for _, u := range r.Users {
			// 			u.Email = os.ExpandEnv(u.Email)
			// 			u.Password = os.ExpandEnv(u.Password)
			// 		}
			// 		continue
			// 	}
			// 	u, ok := v.(User)
			// 	if ok {
			// 		fmt.Printf("\nUser type - %s", u.Email)
			// 		u.Email = os.ExpandEnv(u.Email)
			// 		u.Password = os.ExpandEnv(u.Password)
			// 		continue
			// 	}
			// 	fmt.Printf("\nNone type")
			// }

			viper.Set(key, value)

		default:
			// fmt.Println("default" + key)
			viper.Set(key, value)
		}

	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct %s", err)
	}

	return config, nil
}
