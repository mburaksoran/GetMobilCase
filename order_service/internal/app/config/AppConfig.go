package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// InitFromConfigFile is to populate AppConfig with values from file
func InitFromConfigFile() (*AppConfig, error) {
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

// AppConfig is a structure, instance of which will be populated from provided yaml config
type AppConfig struct {
	AwsConfig struct {
		Key    string `yaml:"key"`
		Region string `yaml:"region"`
		Secret string `yaml:"secret"`
	} `yaml:"aws_config"`
	SqsHost               string `yaml:"sqs_host"`
	SqsMaxWorkerCount     int    `yaml:"sqs_max_worker_count"`
	SqsMaxMessageCount    int    `yaml:"sqs_max_message_count"`
	SqlDatabaseName       string `yaml:"sql_database_name"`
	SqlHost               string `yaml:"sql_host"`
	SqlPassword           string `yaml:"sql_password"`
	SqlPort               string `yaml:"sql_port"`
	SqlUser               string `yaml:"sql_user"`
	SqlSslMode            string `yaml:"sql_ssl_mode"`
	Environment           string `yaml:"environment"`
	MongoDbUrl            string `yaml:"mongo_db_url"`
	MongoDbName           string `yaml:"mongo_db_name"`
	MongoDbCollectionName string `yaml:"mongo_db_collection_name"`
}
