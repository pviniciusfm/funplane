package v1alpha1

// Route defines a route destination for gateway
type Route struct {
	Match         string          `json:"match,omitempty"`
	RouteType     RouteType       `json:"type,omitempty"`
	Service       *ServiceEntry   `json:"service,omitempty"`
	Host          string          `json:"host,omitempty"`
	Rewrite       *PreffixRewrite `json:"rewrite,omitempty"`
	TrafficPolicy *TrafficPolicy  `json:"trafficPolicy,omitempty"`
}

//go:generate jsonenums -type=RouteType
type RouteType int
const (
	CONSUL_SDS RouteType = iota
	CONSUL_DNS
	DNS
)

// PreffixRewrite can be used to rewrite specific parts of a HTTP request
// before forwarding the request to the destination.
type PreffixRewrite struct {
	URI string `json:"uri"`
}

// ServiceEntry defines which service will be used for envoy upstream
type ServiceEntry struct {
	Name string `json:"name,omitempty"`
	Tag  string `json:"tag,omitempty"`
}
