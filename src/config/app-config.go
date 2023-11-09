package config

import "github.com/Netflix/go-env"

var appConfig *AppConfiguration

type AppConfiguration struct {
	HTTPPort uint16 `env:"HTTP_PORT,default=8000"`

	DBHost     string `env:"DB_HOST,default=localhost"`
	DBPort     uint16 `env:"DB_PORT,default=26257"`
	DBUser     string `env:"DB_USER,default=donation_mgmt_app"`
	DBPassword string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_NAME,default=donationsdb"`
	DBSchema   string `env:"DB_SCHEMA,default=donations"`
}

func Bootstrap() *AppConfiguration {
	appConfig = &AppConfiguration{}
	_, err := env.UnmarshalFromEnviron(appConfig)
	if err != nil {
		panic("Error reading environment variables: " + err.Error())
	}

	return appConfig
}

func AppConfig() *AppConfiguration {
	if appConfig == nil {
		panic("AppConfig was not bootstrapped")
	}

	return appConfig
}
