package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GoogleConfig(config *AppConfig) oauth2.Config {
	config.GoogleConfig = oauth2.Config{
		RedirectURL:  config.GoogleRedirectURL,
		ClientID:     config.GoogleClientID,
		ClientSecret: config.GoogleClientSecret,
		Scopes:       config.GoogleScopes,
		Endpoint:     google.Endpoint,
	}

	return config.GoogleConfig
}
