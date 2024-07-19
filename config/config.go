package config

import (
	"elvy-loadbalancer/loadbalancer"
	"fmt"
	ioutil "io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Type     string                 `yaml:"type"`
	Port     string                 `yaml:"port"`
	Backends []loadbalancer.Backend `yaml:"backends"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	return &config, nil
}
