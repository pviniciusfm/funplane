package v1alpha1

//go:generate jsonenums -type=LoadBalancerType
type LoadBalancerType int
const (
	ROUND_ROBIN LoadBalancerType = iota
	LEAST_CONN
	RANDOM
	PASSTHROUGH
)

// LoadBalancer defines a strategy for node load balancer
type LoadBalancerSettings struct {
	LoadBalancerType LoadBalancerType `json:"type,omitempty"`
}
