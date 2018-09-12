package model

// "github.frg.tech/cloud/fanplane/pkg/model/networking/core/v1alpha3"

// ConfigGenerator represents the interfaces to be implemented by code that generates xDS responses
type ConfigGenerator interface {
	// BuildSharedPushState precomputes shared state across all envoys such as outbound clusters, outbound listeners for sidecars,
	// routes for sidecars/gateways etc. This state is stored in the ConfigGenerator object, and reused during the
	// actual BuildClusters/BuildListeners/BuildHTTPRoutes calls. The plugins will not be invoked during this call
	// as most plugins require information about the specific proxy (e.g., mixer/authn/authz). Instead, the
	// BuildYYY functions will invoke the plugins on the precomputed objects.
	// BuildSharedPushState(env *model.Environment, push *model.PushContext) error

	// BuildListeners returns the list of inbound/outbound listeners for the given proxy. This is the LDS output
	// Internally, the computation will be optimized to ensure that listeners are computed only
	// once and shared across multiple invocations of this function.
	// BuildListeners() ([]*v2.Listener, error)

	// BuildClusters returns the list of clusters for the given proxy. This is the CDS output
	// BuildClusters() ([]*v2.Cluster, error)

	// BuildHTTPRoutes returns the list of HTTP routes for the given proxy. This is the RDS output
	// BuildHTTPRoutes(routeName string) (*v2.RouteConfiguration, error)
}

// NewConfigGenerator creates a new instance of the dataplane configuration generator
func NewConfigGenerator(plugins []string) ConfigGenerator {
	return nil
}
