package config

import (
	"fmt"
	"log/slog"

	"github.com/Netflix/go-env"
)

var appConfig *AppConfiguration

type AppEnvironment string

const (
	Development     AppEnvironment = "development"
	IntegrationTest AppEnvironment = "int-tests"
	Staging         AppEnvironment = "staging"
	Production      AppEnvironment = "production"
)

type HTTPAuthenticationMethod string

const (
	AuthFirebase  HTTPAuthenticationMethod = "firebase"
	AuthDevHeader HTTPAuthenticationMethod = "devheader"
)

type AppConfiguration struct {
	HTTPPort uint16 `env:"HTTP_PORT,default=8000"`

	AppName        string         `env:"APP_NAME,default=donation-mgmt"`
	AppEnvironment AppEnvironment `env:"APP_ENVIRONMENT,default=development"`

	LogLevel     string `env:"LOG_LEVEL,default=info"`
	LogAddSource bool   `env:"LOG_ADD_SOURCE,default=false"`

	DBHost     string `env:"DB_HOST,default=localhost"`
	DBPort     uint16 `env:"DB_PORT,default=26257"`
	DBUser     string `env:"DB_USER,default=donation_mgmt_app"`
	DBPassword string `env:"DB_PASSWORD,required=true"`
	DBName     string `env:"DB_NAME,default=donationsdb"`
	DBSchema   string `env:"DB_SCHEMA,default=donations"`

	RewriteForbiddenErrors bool `env:"REWRITE_FORBIDDEN,default=true"`

	HTTPAuthenticationMethod HTTPAuthenticationMethod `env:"HTTP_AUTH,default=firebase"`

	GoogleProjectID           string `env:"GOOGLE_PROJECT_ID"`
	GCPServiceAccountJSONPath string `env:"GCP_SA_JSON_PATH"`
}

func Bootstrap() *AppConfiguration {
	appConfig = &AppConfiguration{}
	_, err := env.UnmarshalFromEnviron(appConfig)
	if err != nil {
		panic("Error reading environment variables: " + err.Error())
	}

	return appConfig
}

func (appConfig *AppConfiguration) WarnUnsafeOptions(logger *slog.Logger) {
	l := logger.With(slog.String("component", "config"))

	if appConfig.AppEnvironment == Development {
		l.Warn(fmt.Sprintf("APP_ENVIRONMENT is set to '%s'. This is unsafe for production environments", appConfig.AppEnvironment))
	}

	if !appConfig.RewriteForbiddenErrors {
		l.Warn("REWRITE_FORBIDDEN is set to disabled. This is unsafe for production environments")
	}

	if appConfig.HTTPAuthenticationMethod == AuthDevHeader {
		l.Warn(fmt.Sprintf("HTTP_AUTH is set to '%s'. This is unsafe for production environments", appConfig.HTTPAuthenticationMethod))
	}

	if appConfig.GCPServiceAccountJSONPath != "" {
		l.Warn("GCP services are authenticated through a service account instead of Google Application Default Credentials. This is not recommended for production environments")
	}
}

func (appConfig *AppConfiguration) EnableFirebase() bool {
	return appConfig.HTTPAuthenticationMethod == AuthFirebase
}

func AppConfig() *AppConfiguration {
	if appConfig == nil {
		panic("AppConfig was not bootstrapped")
	}

	return appConfig
}
