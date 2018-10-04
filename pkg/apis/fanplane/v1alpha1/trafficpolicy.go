package v1alpha1

const (
	DefaultMaxPendingRequest = 1024
	DefaultMaxRequests = 1024
)

// ConnectionPoolSettings defines the settings for envoy's connection pool.
// The settings apply to each individual host in the upstream service. See Envoy’s circuit breaker
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/circuit_breaking
type ConnectionPoolSettings struct {
	//Maximum number of pending HTTP requests to a destination. Default 1024.
	Http1MaxPendingRequests int32 `json:"httpMaxPendingRequests"`
	//Maximum number of requests to a backend. Default 1024.
	Http2MaxRequests int32 `json:"http2MaxRequests"`
	//Maximum number of requests per connection to a backend. Setting this parameter to 1 disables keep alive.
	MaxRequestsPerConnection int32 `json:"maxRequestsPerConnection"`
	//Connection timeout
	Timeout *ReadableDuration `json:"timeout"`
}

// RetryPolicy Describes the retry policy to use when a HTTP request fails
type RetryPolicy struct {
	RetryOn       []string          `json:"retryOn,omitempty"`
	NumRetries    int32             `json:"numRetries,omitempty"`
	PerTryTimeout *ReadableDuration `json:"perTryTimeout,omitempty"`
	//MaxTimeout specifies an overall request timeout.
	//This avoids long request times due to a large number of retries. Default 15s
	MaxTimeout    *ReadableDuration `json:"maxTimeout,omitempty"`
}
