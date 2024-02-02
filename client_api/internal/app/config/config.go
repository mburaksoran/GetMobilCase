package config

import (
	"fmt"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type AppConfig struct {
	GoogleConfig       oauth2.Config
	AwsConfig          AWSConfig `yaml:"aws_config"`
	SqsHost            string    `yaml:"sqs_host"`
	SqlDatabaseName    string    `yaml:"sql_database_name"`
	SqlHost            string    `yaml:"sql_host"`
	SqlPassword        string    `yaml:"sql_password"`
	SqlPort            string    `yaml:"sql_port"`
	SqlUser            string    `yaml:"sql_user"`
	SqlSslMode         string    `yaml:"sql_ssl_mode"`
	Environment        string    `yaml:"environment"`
	MongoDbUrl         string    `yaml:"mongo_db_url"`
	GoogleRedirectURL  string    `yaml:"google_redirect_url"`
	GoogleClientID     string    `yaml:"google_client_id"`
	GoogleClientSecret string    `yaml:"google_client_secret"`
	GoogleScopes       []string  `yaml:"google_scopes"`
}

// InitConfig is to populate AppConfig with values from file
func InitConfig() (*AppConfig, error) {
	data, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		return nil, err
	}

	var config *AppConfig

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to load AppConfig %w", err)
	}

	return config, nil
}
