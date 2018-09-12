package model

import (
	"fmt"
	"io/ioutil"

	validator "gopkg.in/validator.v2"
	yaml "gopkg.in/yaml.v2"
)

type Listener struct {
	Port         uint
	ProtocolType string
}

type Route struct {
}

type Server struct {
}

type Service struct {
}

type TrafficPolicy struct {
}

type Gateway struct {
	Name      string            `yaml:"name"`
	Namespace string            `yaml:"namespace"`
	Subset    map[string]string `yaml:subset`
}

// LoadGateway reads configuration data from a YAML file
func LoadGateway(path string) (*Gateway, error) {
	cfgBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %s", err)
	}
	cfg := new(Gateway)
	err = yaml.Unmarshal(cfgBytes, cfg)
	if err != nil {
		return nil, fmt.Errorf("parsing config: %s", err)
	}
	err = validator.Validate(cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid config: %s", err)
	}
	return cfg, nil
}

// Save writes configuration data to a YAML file
func (c *Gateway) Save(path string) error {
	configBytes, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, configBytes, 0600)
}
