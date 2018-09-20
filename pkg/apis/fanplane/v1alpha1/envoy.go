package v1alpha1

import (
	envoy "github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EnvoyBootstrap is a specification for an raw envoy configuration
type EnvoyBootstrap struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Status GatewayStatus     `json:"status"`
	Tags   map[string]string `json:"tags"`

	Spec map[string]interface{} `json:"spec"`
}


// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// EnvoyList is a list of Raw Envoy resources
type EnvoyBootstrapList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []EnvoyBootstrap `json:"items"`
}

func (in *EnvoyBootstrap) DeepCopyInto(out *EnvoyBootstrap) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = in.Spec
	return
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