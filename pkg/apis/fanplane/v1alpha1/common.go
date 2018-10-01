package v1alpha1

import (
	"bufio"
	"errors"
	"fmt"
	"github.frg.tech/cloud/fanplane/pkg/apis/fanplane"
	"gopkg.in/validator.v2"
	"k8s.io/apimachinery/pkg/util/yaml"
	"os"
)

const (
	MsgConvertError           = "could not convert file %q to fanplane object %q"
	EnvoyBootstrapType string = "EnvoyBootstrap"
	GatewayType               = "Gateway"
)

// IsValidFanplaneObject returns if the kind matches to any entity name in Fanplane models
func IsValidFanplaneObject(kind string) bool {
	switch kind {
	case EnvoyBootstrapType, GatewayType:
		return true
	}

	return false
}

// Status is the status for a Fanplane resource
type Status struct {
	Processed bool `json:"processed"`
}

// LoadFanplaneKind parses only the CRD metadata
func LoadFanplaneKind(path string) (*fanplane.Kind, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Couldn't reader %s filestore", path))
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	cfg := &fanplane.Kind{}
	decoder := yaml.NewYAMLOrJSONDecoder(reader, reader.Size())

	if err := decoder.Decode(cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %s", err)
	}

	if err := validator.Validate(cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %s", err)
	}

	return cfg, nil
}
