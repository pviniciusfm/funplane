package adapter

import (
	"fmt"
	envoy "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
	envoyCluster "github.com/envoyproxy/go-control-plane/envoy/api/v2/cluster"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	envoyRoute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	hcm "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/envoyproxy/go-control-plane/pkg/util"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.frg.tech/cloud/fanplane/pkg/apis/fanplane/v1alpha1"
)

const (
	DefaultListenerName    = "listener_0"
	DefaultStatPrefix      = "ingress_http"
	DefaultVirtualHostName = "local_service"
	DefaultRouteConfigName = "local_route"
)

type GatewayAdapter struct {
	Aggregated
}

type ConversionType int

const (
	LDS = iota
	CDS
	EDS
	RDS
)

var (
	ListenerConversionMap map[v1alpha1.RouteType]func(in *v1alpha1.Route) (*envoyRoute.Route, error)
	ClusterConversionMap  map[v1alpha1.RouteType]func(in *v1alpha1.Route) (*envoy.Cluster, error)
)

func init() {
	ListenerConversionMap = map[v1alpha1.RouteType]func(in *v1alpha1.Route) (*envoyRoute.Route, error){}
	ClusterConversionMap = map[v1alpha1.RouteType]func(in *v1alpha1.Route) (*envoy.Cluster, error){}

	//Route Conversion maps based on RouteType
	ListenerConversionMap[v1alpha1.CONSUL_SDS] = buildConsulSDSRoute
	ListenerConversionMap[v1alpha1.CONSUL_DNS] = buildConsulDNSRoute
	ListenerConversionMap[v1alpha1.DNS] = buildDNSRoute

	//Cluster Conversion maps based on RouteType
	ClusterConversionMap[v1alpha1.CONSUL_SDS] = buildConsulSDSCluster
	ClusterConversionMap[v1alpha1.CONSUL_DNS] = buildConsulDNSCluster
	ClusterConversionMap[v1alpha1.DNS] = buildDNSCluster
}

//GetListeners extract and adapts the single gateway listener in Gateway Fanplane to envoy object
func (GatewayAdapter) GetListeners(object v1alpha1.FanplaneObject) ([]cache.Resource, error) {
	gateway := object.(*v1alpha1.Gateway)
	lis, err := buildListener(gateway)

	if err != nil {
		return nil, err
	}

	return []cache.Resource{lis}, nil
}

//GetRoutes is not implemented. Gateway supports only listeners and clusters from aggregated interface
func (GatewayAdapter) GetRoutes(object v1alpha1.FanplaneObject) (result []cache.Resource, err error) {
	return
}

//GetEndpoints is not implemented. Gateway supports only listeners and clusters from aggregated interface
func (GatewayAdapter) GetEndpoints(object v1alpha1.FanplaneObject) (result []cache.Resource, err error) {
	return
}

//GetClusters extract and adapts Gateway Fanplane endpoints to envoy object
func (GatewayAdapter) GetClusters(object v1alpha1.FanplaneObject) (result []cache.Resource, err error) {
	gateway := object.(*v1alpha1.Gateway)
	result = []cache.Resource{}

	for _, route := range gateway.Spec.Routes {
		cluster, err := buildEnvoyCluster(route)
		if err != nil {
			return nil, err
		}
		result = append(result, cluster)
	}

	return
}

//buildListener produces the required envoy listeners to support Fanplane highlevel Gateway object
func buildListener(in *v1alpha1.Gateway) (out *envoy.Listener, err error) {

	manager, err := buildConnectionManager(in)
	config, err := util.MessageToStruct(manager)

	if err != nil {
		return nil, err
	}

	out = &envoy.Listener{
		Name:    DefaultListenerName,
		Address: BuildAddress(v1alpha1.LocalhostAddress, uint32(in.Spec.Listener.Port)),
		FilterChains: []listener.FilterChain{{
			Filters: []listener.Filter{{
				Name:   util.HTTPConnectionManager,
				Config: config,
			}},
		}},
	}

	return
}

//buildConnectionManager produces the  HttpConnectionManager filter required in envoy listeners
func buildConnectionManager(in *v1alpha1.Gateway) (manager *v2.HttpConnectionManager, err error) {
	vHost, err := buildVirtualHosts(in)

	if err != nil {
		return nil, errors.Errorf("Couldn't build connection manager. %s", err)
	}

	manager = &hcm.HttpConnectionManager{
		CodecType:  hcm.AUTO,
		StatPrefix: DefaultStatPrefix,
		RouteSpecifier: &hcm.HttpConnectionManager_RouteConfig{
			RouteConfig: &envoy.RouteConfiguration{
				Name:         DefaultRouteConfigName,
				VirtualHosts: []envoyRoute.VirtualHost{*vHost},
			},
		},
		HttpFilters: []*hcm.HttpFilter{{
			Name: util.Router,
		}},
	}

	return
}

//buildVirtualHosts extracts all virtual hosts from Gateway routes
func buildVirtualHosts(gateway *v1alpha1.Gateway) (virtualHost *envoyRoute.VirtualHost, err error) {
	routes := gateway.Spec.Routes
	transformedRoutes := make([]envoyRoute.Route, len(routes))

	for idx, route := range routes {
		value, err := buildEnvoyRoute(route)
		if err != nil {
			return nil, fmt.Errorf("couldn't build virtual host for gateway %s", gateway.Name)
		}
		transformedRoutes[idx] = *value
	}

	virtualHost = &envoyRoute.VirtualHost{
		Name:    DefaultVirtualHostName,
		Domains: gateway.Spec.Hosts,
		Routes:  transformedRoutes,
	}

	if gateway.Spec.Hosts == nil {
		virtualHost.Domains = []string{"*"}
	}

	return
}

//buildEnvoyRoute gets the function to translate routes accordingly with RouteType
func buildEnvoyRoute(in *v1alpha1.Route) (route *envoyRoute.Route, err error) {
	convert := ListenerConversionMap[in.RouteType]
	route, err = convert(in)
	return
}

//buildEnvoyRoute gets the function to translate clusters accordingly with RouteType
func buildEnvoyCluster(in *v1alpha1.Route) (cluster *envoy.Cluster, err error) {
	convert := ClusterConversionMap[in.RouteType]
	cluster, err = convert(in)
	return
}

//buildConsulDNSRoute produces an envoy RDS Route from Gateway RouteType
func buildConsulDNSRoute(in *v1alpha1.Route) (out *envoyRoute.Route, err error) {
	//TODO: See if we need to support HostRewrite
	//HostRewriteSpecifier:

	out = &envoyRoute.Route{
		Match: envoyRoute.RouteMatch{
			PathSpecifier: &envoyRoute.RouteMatch_Prefix{
				Prefix: in.Prefix,
			},
		},
		Action: &envoyRoute.Route_Route{
			Route: &envoyRoute.RouteAction{
				ClusterSpecifier: &envoyRoute.RouteAction_Cluster{
					Cluster: in.Service.GetClusterId(),
				},
			},
		},
	}

	route := out.GetAction().(*envoyRoute.Route_Route).Route
	if in.Rewrite != nil && in.Rewrite.URI != "" {
		route.HostRewriteSpecifier = &envoyRoute.RouteAction_HostRewrite{
			HostRewrite: in.Rewrite.URI,
		}
	} else {
		route.HostRewriteSpecifier = &envoyRoute.RouteAction_AutoHostRewrite{
			AutoHostRewrite: &types.BoolValue{Value: true},
		}
	}

	return
}

//buildConsulSDSRoute produces envoy endpoints from consul catalog api calls
//WARNING: Not implemented this will be supported in future versions
func buildConsulSDSRoute(in *v1alpha1.Route) (out *envoyRoute.Route, err error) {
	//NOT SUPPORTED IN V1ALPHA1
	return
}

//buildDNSRoute produces envoy endpoints from DNS gateway RouteType
func buildDNSRoute(in *v1alpha1.Route) (out *envoyRoute.Route, err error) {

	action := &envoyRoute.RouteAction{
		ClusterSpecifier: &envoyRoute.RouteAction_Cluster{
			Cluster: in.DNS.GetClusterId(),
		},
		HostRewriteSpecifier: &envoyRoute.RouteAction_AutoHostRewrite{
			&types.BoolValue{Value: true},
		},
	}

	if in.Rewrite != nil && in.Rewrite.URI != "" {
		action.PrefixRewrite = in.Rewrite.URI
	}

	buildRetryPolicy(in, action)
	if in.ConnectionPoolSettings != nil {
		if in.ConnectionPoolSettings.Timeout.Duration != 0 {
			action.Timeout = &in.ConnectionPoolSettings.Timeout.Duration
		}
	}

	out = &envoyRoute.Route{
		Match: envoyRoute.RouteMatch{
			PathSpecifier: &envoyRoute.RouteMatch_Prefix{
				Prefix: in.Prefix,
			},
		},
		Action: &envoyRoute.Route_Route{
			Route: action,
		},
	}

	return
}

//buildRetryPolicy translates Gateway retryPolicy objects to envoy configuration
func buildRetryPolicy(in *v1alpha1.Route, action *envoyRoute.RouteAction) {
	if in.RetryPolicy != nil {
		action.RetryPolicy = &envoyRoute.RouteAction_RetryPolicy{}
		numRetries := uint32(in.RetryPolicy.NumRetries)

		if numRetries != 0 {
			action.RetryPolicy.NumRetries = &types.UInt32Value{Value: uint32(numRetries)}
		}
		perTryTimeout := in.RetryPolicy.PerTryTimeout.Duration

		if perTryTimeout != 0 {
			action.RetryPolicy.PerTryTimeout = &perTryTimeout
		}
	}
}

//buildConsulSDSCluster extracts envoy's cluster definition from consul catalog api
//WARNING: Not implemented this will be supported in future versions
func buildConsulSDSCluster(in *v1alpha1.Route) (out *envoy.Cluster, err error) {
	//NOT SUPPORTED IN V1ALPHA1
	return
}

func buildConsulDNSCluster(in *v1alpha1.Route) (out *envoy.Cluster, err error) {
	domain := viper.GetString("domain")
	host := buildConsulAddress(in.Service.Tag, in.Service.Name, in.Service.Datacenter, domain)
	address := BuildAddress(host, in.Service.Port)

	out = &envoy.Cluster{
		Name:            in.Service.GetClusterId(),
		Hosts:           []*core.Address{&address},
		Type:            envoy.Cluster_STRICT_DNS,
		DnsLookupFamily: envoy.Cluster_V4_ONLY,
		LbPolicy:        envoy.Cluster_LbPolicy(in.GetLoadBalancerSettings().LoadBalancerType),
	}

	if in.ConnectionPoolSettings != nil && in.ConnectionPoolSettings.Timeout.Duration != 0 {
		out.ConnectTimeout = in.ConnectionPoolSettings.Timeout.Duration
	}

	buildConnectionPoolSettings(in, out)

	if in.TLSContext != nil && in.TLSContext.SNI != "" {
		out.TlsContext = &auth.UpstreamTlsContext{
			Sni: in.TLSContext.SNI,
		}
	}

	return
}

func buildConnectionPoolSettings(in *v1alpha1.Route, out *envoy.Cluster) {
	if in.ConnectionPoolSettings != nil {
		thresholds := &envoyCluster.CircuitBreakers_Thresholds{
			MaxConnections:     &types.UInt32Value{Value: uint32(in.ConnectionPoolSettings.MaxRequestsPerConnection)},
			MaxPendingRequests: &types.UInt32Value{Value: uint32(in.ConnectionPoolSettings.Http1MaxPendingRequests)},
		}

		out.CircuitBreakers = &envoyCluster.CircuitBreakers{
			Thresholds: []*envoyCluster.CircuitBreakers_Thresholds{thresholds},
		}
	}
}

func buildDNSCluster(in *v1alpha1.Route) (out *envoy.Cluster, err error) {

	address := BuildAddress(in.DNS.Address, uint32(in.DNS.Port))
	out = &envoy.Cluster{
		Name:            in.DNS.GetClusterId(),
		Hosts:           []*core.Address{&address},
		Type:            envoy.Cluster_DiscoveryType(in.DNS.Type),
		DnsLookupFamily: envoy.Cluster_V4_ONLY,
		LbPolicy:        envoy.Cluster_LbPolicy(in.GetLoadBalancerSettings().LoadBalancerType),
	}

	if in.ConnectionPoolSettings != nil && in.ConnectionPoolSettings.Timeout.Duration != 0 {
		out.ConnectTimeout = in.ConnectionPoolSettings.Timeout.Duration
	}

	if in.TLSContext != nil && in.TLSContext.SNI != "" {
		out.TlsContext = &auth.UpstreamTlsContext{
			Sni: in.TLSContext.SNI,
		}
	}

	return
}

func buildConsulAddress(tag string, serviceName string, datacenter string, domain string) (dns string) {
	dns = fmt.Sprintf("%s.service", serviceName)

	if tag != "" {
		dns = fmt.Sprintf("%s.%s", tag, dns)
	}

	if datacenter != "" {
		dns = fmt.Sprintf("%s.%s", dns, datacenter)
	}

	dns = fmt.Sprintf("%s.%s", dns, domain)

	return
}
