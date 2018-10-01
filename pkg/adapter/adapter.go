//adapter package is responsible for the Translation of Fanplane Objects to the XDS Envoy Models
package adapter

import (
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.frg.tech/cloud/fanplane/pkg/apis/fanplane"
	"github.frg.tech/cloud/fanplane/pkg/apis/fanplane/v1alpha1"
)

type Adapter interface {
	GetFanplaneObject() fanplane.Object
	XdsResponseBuilder
}

//New adapter is a factory method to return the concrete adapter
func NewAdapter(obj fanplane.Object) Adapter {
	switch obj.GetObjectKind().GroupVersionKind().Kind {
	case v1alpha1.EnvoyBootstrapType:
		return NewEnvoyBootstrapAdapter(obj)
	case v1alpha1.GatewayType:
		return NewGatewayAdapter(obj)
	default:
		return nil
	}
}

// XdsResponseBuilder represents the interfaces to be implemented by code that generates xDS responses
type XdsResponseBuilder interface {
	// BuildListeners returns the list of inbound/outbound listeners for the given proxy. This is the LDS output
	BuildListeners() ([]cache.Resource, error)
	// BuildRoutes returns the list of HTTP routes for the given proxy. This is the RDS output
	BuildRoutes() ([]cache.Resource, error)
	// BuildClusters returns the list of clusters for the given proxy. This is the CDS output
	BuildClusters() ([]cache.Resource, error)
}
