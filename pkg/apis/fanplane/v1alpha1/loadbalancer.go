package v1alpha1

//go:generate jsonenums -type=LoadBalancerType
type LoadBalancerType int32

const (
	ROUND_ROBIN LoadBalancerType = iota
	LEAST_CONN
	_
	RANDOM
)

// LoadBalancer defines a strategy for node load balancer
type LoadBalancerSettings struct {
	LoadBalancerType LoadBalancerType `json:"type,omitempty"`
}
