package config

import "github.com/Netflix/go-env"

var appConfig *AppConfiguration

type AppConfiguration struct {
	HttpPort uint16 `env:"HTTP_PORT,default=8000"`
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
