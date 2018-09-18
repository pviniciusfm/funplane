package model

import (
	"fmt"
	"io/ioutil"

	validator "gopkg.in/validator.v2"
	yaml "gopkg.in/yaml.v2"
)

//Listener defines envoy tcp listeners
type Listener struct {
	Port         uint
	ProtocolType string
}

//Route define a route in fanplane
type Route struct {
}

//Server defines the grpc server for fanplane
type Server struct {
}

//Service Defines a service in fanplane
type Service struct {
}

//TrafficPolicy defines how the load will be distributed
type TrafficPolicy struct {
}

//Gateway defines an egress gateway
type Gateway struct {
	Name      string            `yaml:"name"`
	Namespace string            `yaml:"namespace"`
	Subset    map[string]string `yaml:"subset"`
}

// LoadGateway reads configuration data from a YAML file
func LoadGateway(path string) (cfg *Gateway, err error) {
	cfgBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %s", err)
	}

	cfg = &Gateway{}
	err = yaml.Unmarshal(cfgBytes, cfg)

	if err != nil {
		return nil, fmt.Errorf("parsing config: %s", err)
	}

	err = validator.Validate(cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid config: %s", err)
	}

	return
}

// Save writes configuration data to a YAML file
func (c *Gateway) Save(path string) error {
	configBytes, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, configBytes, 0600)
}
