package v1alpha1

import (
	"k8s.io/apimachinery/pkg/util/uuid"
)

// Route defines a route destination for gateway
type Route struct {
	clusterId  string
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
	FaultInjection         *FaultInjection         `json:"faultInjection,omitempty"`
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

func (in *Route) GetClusterId() string {
	if in.clusterId == "" {
		in.clusterId = string(uuid.NewUUID())
	}
	return in.clusterId
}

//DNSRoute is used to represent a Logical DNS Route Type Entry
type DNSRoute struct {
	Type    *DNSType `json:"type,omitempty"`
	Address string   `json:"address,omitempty"`
	Port    int32    `json:"port,omitempty"`
}

//go:generate jsonenums -type=DNSType
type DNSType int

const (
	STATIC DNSType = iota
	STRICT_DNS
	LOGICAL_DNS
	EDS
)

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

//FaultInjection makes possible experiments of connection reliability
type FaultInjection struct {
	//AbortPercent the percentage of operations/connection requests on which the envoy will abort connections.
	AbortPercent *float32 `json:"abortPercent,omitempty"`
	//Delay add a fixed delay before forwarding the operation upstream
	Delay *ReadableDuration `json:"fixedDelay,omitempty"`
	//DelayPercent the percentage of operations/connection requests on which the delay will be injected.
	DelayPercent *float32 `json:"delayPercent,omitempty"`
}
