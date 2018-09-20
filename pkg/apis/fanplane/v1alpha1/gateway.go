package v1alpha1

import (
	"bufio"
	"errors"
	"fmt"
	"gopkg.in/validator.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"os"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Gateway is a specification for an Gateway resource
type Gateway struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Status GatewayStatus `json:"status,omitempty"`

	Spec *GatewaySpec `json:"spec" validate:"nonzero"`
}

type GatewaySpec struct {
	Selector map[string]string `json:"selector,omitempty"`
	Hosts    []string          `json:"hosts"`
	Listener *Listener         `json:"listener,omitempty" validate:"nonzero"`
	Routes   []*Route          `json:"routes"`
}

// GatewayStatus is the status for a Gateway resource
type GatewayStatus struct {
	Processed bool `json:"processed"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GatewayList is a list of Gateway resources
type GatewayList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Gateway `json:"items"`
}

// LoadGateway reads configuration data from a YAML file
func LoadGateway(path string) (cfg *Gateway, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Couldn't reader %s file", path))
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	cfg = &Gateway{}
	decoder := yaml.NewYAMLOrJSONDecoder(reader, reader.Size())
	err = decoder.Decode(cfg)

	if err != nil {
		return nil, fmt.Errorf("parsing config: %s", err)
	}

	err = validator.Validate(cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid config: %s", err)
	}

	return
}
