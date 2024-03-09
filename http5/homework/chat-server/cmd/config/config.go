package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	StorageType string   `yaml:"storage_type"`
	AuthType    string   `yaml:"auth_type"`
	Database    DBConfig `yaml:"database"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DB       string `yaml:"db"`
}

func (c *Config) ParseData(path string) error {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return err
	}

	return nil
}
