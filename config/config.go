package config

import (
	"os"

	"github.com/ZooLearn/file/internal/rabbitx"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
)

func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

type Config struct {
	EnvConf      EnvConf                 `yaml:"EnvConf"`
	ProducerConf rabbitx.RabbitmqProConf `yaml:"ProducerConf"`
	ConsumerConf rabbitx.RabbitmqConConf `yaml:"ConsumerConf"`
}

type EnvConf struct {
	AppEnv                 string `yaml:"AppEnv"`
	ServerAddress          string `yaml:"ServerAddress"`
	ContextTimeout         int    `yaml:"ContextTimeout"`
	AccessTokenExpiryHour  int    `yaml:"AccessTokenExpiryHour"`
	RefreshTokenExpiryHour int    `yaml:"RefreshTokenExpiryHour"`
	AccessTokenSecret      string `yaml:"AccessTokenSecret"`
	RefreshTokenSecret     string `yaml:"RefreshTokenSecret"`
}
