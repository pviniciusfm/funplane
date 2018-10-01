package adapter

import (
	"fmt"
	envoy "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
	envoyCluster "github.com/envoyproxy/go-control-plane/envoy/api/v2/cluster"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	envoyRoute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	xdsfault "github.com/envoyproxy/go-control-plane/envoy/config/filter/fault/v2"
	xdshttpfault "github.com/envoyproxy/go-control-plane/envoy/config/filter/http/fault/v2"
	"github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	hcm "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/envoyproxy/go-control-plane/envoy/type"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/envoyproxy/go-control-plane/pkg/util"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.frg.tech/cloud/fanplane/pkg/apis/fanplane"
	"github.frg.tech/cloud/fanplane/pkg/apis/fanplane/v1alpha1"
	"gopkg.in/validator.v2"
	"time"
)

const (
	MsgNotSupported        = "%s is not supported in current version"
	DefaultListenerName    = "listener_0"
	DefaultStatPrefix      = "ingress_http"
	DefaultVirtualHostName = "local_service"
	DefaultRouteConfigName = "local_route"
	DefaultConnectTimeout  = time.Second * 5
)

type gatewayAdapter struct {
	gateway *v1alpha1.Gateway
}

func (adapt *gatewayAdapter) GetFanplaneObject() fanplane.Object {
	return adapt.gateway
}

type ConversionType int

const (
	LDS = iota
	CDS
	EDS
	RDS
)

var ClusterConversionMap map[v1alpha1.RouteType]func(in *v1alpha1.Route, out *envoy.Cluster) (error)

func init() {
	ClusterConversionMap = map[v1alpha1.RouteType]func(in *v1alpha1.Route, out *envoy.Cluster) (error){}

	//Cluster Conversion maps based on RouteType
	ClusterConversionMap[v1alpha1.CONSUL_SDS] = buildConsulSDSCluster
	ClusterConversionMap[v1alpha1.CONSUL_DNS] = buildConsulDNSCluster
	ClusterConversionMap[v1alpha1.DNS] = buildRawDNSCluster
}

// NewGatewayAdapter factory method for adapter constructor
func NewGatewayAdapter(obj fanplane.Object) (Adapter) {
	result := &gatewayAdapter{}

	if _, ok := obj.(*fanplane.Kind); ok {
		result.gateway = &v1alpha1.Gateway{TypeMeta: obj.GetTypeMeta(), ObjectMeta: obj.GetObjectMeta()}
	} else {
		result.gateway = obj.(*v1alpha1.Gateway)
	}

	if _, ok := obj.GetSpec().(*v1alpha1.GatewaySpec); !ok {
		spec := &v1alpha1.GatewaySpec{}
		v1alpha1.ToStruct(obj.GetSpec(), spec)
		result.gateway.Spec = spec
	}

	return result
}

//BuildListeners extract and adapts the single gateway listener in Gateway Fanplane to envoy object
func (adapt *gatewayAdapter) BuildListeners() (result []cache.Resource, err error) {
	manager, err := buildConnectionManager(adapt.gateway)
	if err != nil {
		log.Errorf("failed to build connection manager in %q gateway", adapt.gateway.Name)
		return nil, err
	}

	config, err := util.MessageToStruct(manager)
	if err != nil {
		log.Errorf("failed to convert connection manager struct in %q gateway ", adapt.gateway.Name)
		return nil, err
	}

	out := &envoy.Listener{
		Name:    DefaultListenerName,
		Address: BuildAddress(v1alpha1.LocalhostAddress, uint32(adapt.gateway.Spec.Listener.Port)),
		FilterChains: []listener.FilterChain{{
			Filters: []listener.Filter{{
				Name:   util.HTTPConnectionManager,
				Config: config,
			}},
		}},
	}

	return []cache.Resource{out}, nil
}

//BuildRoutes is not implemented. Gateway supports only listeners and clusters from aggregated interface
func (gatewayAdapter) BuildRoutes() (result []cache.Resource, err error) {
	return
}

//BuildEndpoints is not implemented. Gateway supports only listeners and clusters from aggregated interface
func (gatewayAdapter) BuildEndpoints() (result []cache.Resource, err error) {
	return
}

//BuildClusters extract and adapts Gateway Fanplane endpoints to envoy object
func (adapt *gatewayAdapter) BuildClusters() (result []cache.Resource, err error) {
	gateway := adapt.gateway
	result = []cache.Resource{}

	for _, route := range gateway.Spec.Routes {
		cluster, err := buildCluster(route)
		if err != nil {
			return nil, err
		}
		result = append(result, cluster)
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
		value, err := buildRoute(route)
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

	if gateway.Spec.Hosts == nil || len(gateway.Spec.Hosts) == 0 {
		virtualHost.Domains = []string{"*"}
	}

	return
}

//buildRoute produces an envoy RDS Route from Gateway RouteType
func buildRoute(in *v1alpha1.Route) (out *envoyRoute.Route, err error) {
	log.Debug("parsing dns route")

	err = validator.Validate(in)
	if err != nil {
		return nil, err
	}

	out = &envoyRoute.Route{
		Match: envoyRoute.RouteMatch{
			PathSpecifier: &envoyRoute.RouteMatch_Prefix{
				Prefix: in.Prefix,
			},
		},
		Action: &envoyRoute.Route_Route{
			Route: &envoyRoute.RouteAction{
				ClusterSpecifier: &envoyRoute.RouteAction_Cluster{
					Cluster: in.GetClusterId(),
				},
			},
		},
	}

	action := out.Action.(*envoyRoute.Route_Route).Route
	buildRetryPolicy(in, action)
	buildConnectionPoolSettingsRoute(in, action)
	buildFaultInjection(in, out)

	if in.Rewrite != nil && in.Rewrite.URI != "" {
		action.HostRewriteSpecifier = &envoyRoute.RouteAction_HostRewrite{
			HostRewrite: in.Rewrite.URI,
		}
	} else {
		action.HostRewriteSpecifier = &envoyRoute.RouteAction_AutoHostRewrite{
			AutoHostRewrite: &types.BoolValue{Value: true},
		}
	}

	return
}

func buildFaultInjection(in *v1alpha1.Route, out *envoyRoute.Route) {
	if fault := in.FaultInjection; fault != nil {
		if protoMsg, err := util.MessageToStruct(translateFault(fault)); err == nil {
			out.PerFilterConfig["envoy.fault"] = protoMsg
		}
	}
}

//buildConsulSDSRoute produces envoy endpoints from consul catalog api calls
//WARNING: Not implemented this will be supported in future versions
func buildConsulSDSRoute(*v1alpha1.Route, *envoyRoute.Route) (err error) {
	err = fmt.Errorf(MsgNotSupported, "Consul SDS Route")
	log.Error(err)
	return
}

//buildConnectionPoolSettingsRoute build envoy's connection policies.
func buildConnectionPoolSettingsRoute(in *v1alpha1.Route, action *envoyRoute.RouteAction) {
	duration := DefaultConnectTimeout
	if in.ConnectionPoolSettings != nil && in.ConnectionPoolSettings.Timeout != nil {
		duration = in.ConnectionPoolSettings.Timeout.Duration()
	}

	action.Timeout = &duration
}

//buildRetryPolicy translates Gateway retryPolicy objects to envoy configuration
func buildRetryPolicy(in *v1alpha1.Route, action *envoyRoute.RouteAction) {
	log.Debug("parsing retry policy")

	if in.RetryPolicy != nil {
		action.RetryPolicy = &envoyRoute.RouteAction_RetryPolicy{}
		numRetries := uint32(in.RetryPolicy.NumRetries)

		if numRetries != 0 {
			action.RetryPolicy.NumRetries = &types.UInt32Value{Value: uint32(numRetries)}
		}
		perTryTimeout := in.RetryPolicy.PerTryTimeout

		if perTryTimeout != nil {
			duration := perTryTimeout.Duration()
			action.RetryPolicy.PerTryTimeout = &duration
		}

		if in.RetryPolicy.MaxTimeout != nil {
			duration := in.RetryPolicy.MaxTimeout.Duration()
			action.Timeout = &duration
		}
	}
}

//buildConsulSDSCluster extracts envoy's cluster definition from consul catalog api
//WARNING: Not implemented this will be supported in future versions
func buildConsulSDSCluster(*v1alpha1.Route, *envoy.Cluster) (err error) {
	log.Warnf(MsgNotSupported, "Consul SDS Cluster")
	return
}

func buildCluster(in *v1alpha1.Route) (out *envoy.Cluster, err error) {
	//Default cluster definition
	out = &envoy.Cluster{
		Name:            in.GetClusterId(),
		Type:            envoy.Cluster_STRICT_DNS,
		DnsLookupFamily: envoy.Cluster_V4_ONLY,
		LbPolicy:        envoy.Cluster_LbPolicy(in.GetLoadBalancerSettings().LoadBalancerType),
	}

	buildConnectionPoolSettings(in, out)
	buildTLSContext(in, out)

	//Uses parser according to conversion map
	convert := ClusterConversionMap[in.RouteType]
	err = convert(in, out)

	return
}

func buildTLSContext(in *v1alpha1.Route, out *envoy.Cluster) {
	if in.TLSContext != nil && in.TLSContext.SNI != "" {
		out.TlsContext = &auth.UpstreamTlsContext{
			Sni: in.TLSContext.SNI,
		}
	}
}

func buildConsulDNSCluster(in *v1alpha1.Route, out *envoy.Cluster) (err error) {
	if in.Service == nil {
		err = fmt.Errorf("trying building dns cluster with nil service")
		log.Error(err)
		return
	}

	domain := viper.GetString("domain")
	host := buildConsulAddress(in.Service.Tag, in.Service.Name, in.Service.Datacenter, domain)
	address := BuildAddress(host, in.Service.Port)
	out.Hosts = []*core.Address{&address}
	out.Type = envoy.Cluster_STRICT_DNS
	return
}

func buildConnectionPoolSettings(in *v1alpha1.Route, out *envoy.Cluster) {
	setConnectionPoolDefaults(in)
	out.ConnectTimeout = in.ConnectionPoolSettings.Timeout.Duration()

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

func buildRawDNSCluster(in *v1alpha1.Route, out *envoy.Cluster) (err error) {
	address := BuildAddress(in.DNS.Address, uint32(in.DNS.Port))
	out.Hosts = []*core.Address{&address}
	dnsType := envoy.Cluster_STRICT_DNS
	out.LbPolicy = envoy.Cluster_LbPolicy(in.GetLoadBalancerSettings().LoadBalancerType)

	if in.DNS.Type != nil {
		dnsType = envoy.Cluster_DiscoveryType(out.Type)
	}

	out.Type = dnsType

	return
}

//buildConsulAddress builds fqdn for consul based endpoint. e.g {tag}.{service}.{datacenter}.{domain}
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

//setConnectionPoolDefaults initialize routes with connection pool defaults
func setConnectionPoolDefaults(in *v1alpha1.Route) {
	if in.ConnectionPoolSettings == nil {
		in.ConnectionPoolSettings = &v1alpha1.ConnectionPoolSettings{}
	}

	if in.ConnectionPoolSettings.Timeout == nil {
		in.ConnectionPoolSettings.Timeout = v1alpha1.NewReadableDuration(DefaultConnectTimeout)
	}

	if in.ConnectionPoolSettings.MaxRequestsPerConnection == 0 {
		in.ConnectionPoolSettings.MaxRequestsPerConnection = 1024
	}

	if in.ConnectionPoolSettings.Http1MaxPendingRequests == 0 {
		in.ConnectionPoolSettings.Http1MaxPendingRequests = 1024
	}
}

// translateFault translates networking.HTTPFaultInjection into Envoy's HTTPFault
func translateFault(in *v1alpha1.FaultInjection) *xdshttpfault.HTTPFault {
	if in == nil {
		return nil
	}

	out := xdshttpfault.HTTPFault{}
	if in.Delay != nil {
		out.Delay = &xdsfault.FaultDelay{Type: xdsfault.FaultDelay_FIXED}
		out.Delay.Percentage = translatePercentToFractionalPercent(in.DelayPercent)

		delayDuration := in.Delay.Duration()
		out.Delay.FaultDelaySecifier = &xdsfault.FaultDelay_FixedDelay{
			FixedDelay: &delayDuration,
		}

		if in.AbortPercent != nil {
			out.Abort = &xdshttpfault.FaultAbort{}
			out.Abort.Percentage = translatePercentToFractionalPercent(in.AbortPercent)
		}

		if out.Delay == nil && out.Abort == nil {
			return nil
		}
	}

	return &out
}

// translatePercentToFractionalPercent translates an float pointer
// to an envoy.type.FractionalPercent instance.
func translatePercentToFractionalPercent(p *float32) *envoy_type.FractionalPercent {
	return &envoy_type.FractionalPercent{
		Numerator:   uint32(*p * 10000),
		Denominator: envoy_type.FractionalPercent_MILLION,
	}
}
