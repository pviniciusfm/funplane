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

//Gateway is the highlevel CRD object that holds Gateway Specification object
type Gateway struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Status Status `json:"status,omitempty"`

	Spec *GatewaySpec `json:"spec" validate:"nonzero"`
}

func (in *Gateway) GetType() string {
	return in.TypeMeta.Kind
}

func (in *Gateway) GetSpec() interface{} {
	return in.Spec
}

func (in *Gateway) SetSpec(spec interface{}) {
	in.Spec = spec.(*GatewaySpec)
}

func (in *Gateway) GetObjectMeta() metav1.ObjectMeta {
	return in.ObjectMeta
}

func (in *Gateway) SetObjectMeta(objMeta metav1.ObjectMeta) {
	in.ObjectMeta = objMeta
}

// GatewaySpec contains the definition of Envoy's Routes, Listeners and Clusters
type GatewaySpec struct {
	Selector map[string]string `json:"selector,omitempty"`
	Hosts    []string          `json:"hosts"`
	Listener *Listener         `json:"listener,omitempty" validate:"nonzero"`
	Routes   []*Route          `json:"routes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GatewayList is a list of Gateway resources (Required for k8 crds)
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
