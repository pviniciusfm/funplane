package v1alpha1

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/uuid"
)

const UnknownRouteType = "UNKNOWN"

// Route defines a route destination for gateway
type Route struct {
	Prefix     string              `json:"prefix,omitempty"`
	RouteType  RouteType           `json:"type,omitempty"`
	Service    *ConsulServiceEntry `json:"service,omitempty"`
	DNS        *DNSRoute           `json:"dns,omitempty"`
	Host       string              `json:"host,omitempty"`
	Rewrite    *PrefixRewrite      `json:"rewrite,omitempty"`
	TLSContext *TLSContext         `json:"tlsContext,omitempty"`

	LoadBalancerSettings   *LoadBalancerSettings   `json:"loadBalancer,omitempty"`
	ConnectionPoolSettings *ConnectionPoolSettings `json:"connectionPool,omitempty"`
	RetryPolicy            *RetryPolicy            `json:"retryPolicy,omitempty"`
}

//go:generate jsonenums -type=RouteType
type RouteType int

const (
	CONSUL_SDS RouteType = iota
	CONSUL_DNS
	DNS
)

// PrefixRewrite can be used to rewrite specific parts of a HTTP request
// before forwarding the request to the destination.
type PrefixRewrite struct {
	URI string `json:"uri"`
}

// ConsulServiceEntry defines which service will be used for envoy upstream
type ConsulServiceEntry struct {
	Name       string `json:"name,omitempty"`
	Datacenter string `json:"dc,omitempty"`
	Tag        string `json:"tag,omitempty"`
	Port       uint32 `json:"port"`
}

func (in *ConsulServiceEntry) GetClusterId() string {
	return fmt.Sprintf("%s|%s", in.Name, in.Tag)
}

//String returns string representation of RouteType
func (r RouteType) String() string {
	out := UnknownRouteType

	if _, val := _RouteTypeValueToName[r]; val {
		out = _RouteTypeValueToName[r]
	}

	return out
}

//DNSRoute is used to represent a Logical DNS Route Type Entry
type DNSRoute struct {
	id      string  `json:id,omitempty`
	Type    DNSType `json:"type,omitempty"`
	Address string  `json:"address,omitempty"`
	Port    int32   `json:"port,omitempty"`
}

//go:generate jsonenums -type=DNSType
type DNSType int

const (
	STATIC DNSType = iota
	STRICT_DNS
	LOGICAL_DNS
	EDS
)

func (in *DNSRoute) GetClusterId() string {
	if in.id == "" {
		in.id = string(uuid.NewUUID())
	}
	return in.id
}

func (in *Route) GetLoadBalancerSettings() *LoadBalancerSettings {
	if in.LoadBalancerSettings == nil {
		in.LoadBalancerSettings = &LoadBalancerSettings{}
	}
	return in.LoadBalancerSettings
}

//TLSContext defines Server Name Indication of the upstream connection
type TLSContext struct {
	SNI string `json:"sni,omitempty"`
}
