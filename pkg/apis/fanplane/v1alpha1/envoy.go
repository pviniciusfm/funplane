package v1alpha1

import (
	"fmt"
	envoy "github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v2"
	"github.com/ghodss/yaml"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EnvoyBootstrap is a specification for an raw envoy configuration
type EnvoyBootstrap struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Status Status            `json:"status"`
	Tags   map[string]string `json:"tags"`

	Spec interface{} `json:"spec"`
}

func (in *EnvoyBootstrap) GetTypeMeta() metav1.TypeMeta {
	return in.TypeMeta
}

func (in *EnvoyBootstrap) SetSpec(newSpec interface{}) {
	in.Spec = newSpec
}

func (in *EnvoyBootstrap) GetObjectMeta() metav1.ObjectMeta {
	return in.ObjectMeta
}

func (in *EnvoyBootstrap) SetObjectMeta(newMeta metav1.ObjectMeta) {
	in.ObjectMeta = newMeta
}

func (in *EnvoyBootstrap) GetSpec() interface{} {
	return in.Spec
}

func (in *EnvoyBootstrap) DeepCopyInto(out *EnvoyBootstrap) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = in.Spec
	return
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// EnvoyList is a list of Raw Envoy resources
type EnvoyBootstrapList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []EnvoyBootstrap `json:"items"`
}

// ConvertObject converts an EnvoyObject k8s-style object to the
// internal configuration model.
func ParseEnvoyConfig(fanplaneObject interface{}) (config *envoy.Bootstrap, err error) {
	// parse as JSON
	protoRaw, err := FromJSONMap(fanplaneObject, "envoy.config.bootstrap.v2.Bootstrap")
	if err != nil {
		return
	}
	config = protoRaw.(*envoy.Bootstrap)
	return
}

func (in *EnvoyBootstrap) GetSidecarSelector() string {
	return in.Name
}

//LoadEnvoyBootstrap returns EnvoyBootstrap fanplane object from an yaml file into
func LoadEnvoyBootstrap(path string) (bootstrap *EnvoyBootstrap, err error) {
	bootstrap = &EnvoyBootstrap{}
	//Read bytes from filestore
	in, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf(MsgConvertError, path, "EnvoyBootstrap")
	}

	//Unmarshal type struct with envoy unparsed
	if err = yaml.Unmarshal(in, bootstrap); err != nil {
		return nil, fmt.Errorf(MsgConvertError, path, "EnvoyBootstrap")
	}

	//Parse and validates proto message
	envoyConfig, err := ParseEnvoyConfig(bootstrap.Spec)
	if err != nil {
		return nil, fmt.Errorf(MsgConvertError, path, "EnvoyBootstrap")
	}

	bootstrap.Spec = envoyConfig
	return
}
