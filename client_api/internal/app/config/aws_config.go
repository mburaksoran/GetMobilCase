package config

type AWSConfig struct {
	Key    string `yaml:"key"`
	Region string `yaml:"region"`
	Secret string `yaml:"secret"`
}
