//Server is the module responsible to create and manage GRPC xDS server
package server

import (
	"context"
	"fmt"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	xds "github.com/envoyproxy/go-control-plane/pkg/server"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
)

const (
	typePrefix = "type.googleapis.com/envoy.api.v2."
	// Constants used for XDS
	// ClusterType is used for cluster discovery. Typically first request received
	ClusterType = typePrefix + "Cluster"
	// EndpointType is used for EDS and ADS endpoint discovery. Typically second request.
	EndpointType = typePrefix + "ClusterLoadAssignment"
	// ListenerType is sent after clusters and endpoints.
	ListenerType = typePrefix + "Listener"
	// RouteType is sent after listeners.
	RouteType = typePrefix + "RouteConfiguration"
)

type Server struct {
	xDSserver  xds.Server
	Config     *FanplaneConfig
	Cache      cache.SnapshotCache
	Context    context.Context
	callback   *Callback
	startFuncs []func(<-chan struct{}) (error)
	grpcServer *grpc.Server
}

// Hasher returns node ID as an ID
type Hasher struct{}

// ID function
func (h Hasher) ID(node *core.Node) string {
	if node == nil {
		return "unknown"
	}
	return node.Id
}

const grpcMaxConcurrentStreams = 1000000

// RunServer runs gRPC service to publish ADS service
// RunManagementServer starts an xDS server at the given port.
func NewManagementServer(ctx context.Context, config *FanplaneConfig) (server *Server, err error) {
	cache := cache.NewSnapshotCache(config.ADSMode, Hasher{}, log.StandardLogger())
	cb := &Callback{}
	xdsServer := xds.NewServer(cache, cb)

	//Register server modules
	startFuncs := []func(<-chan struct{}) (error){server.HandleAbortSignal}
	grpcServer := grpc.NewServer(grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams))

	server = &Server{
		xdsServer,
		config,
		cache,
		ctx,
		cb,
		startFuncs,
		grpcServer,
	}

	//Register services
	discovery.RegisterAggregatedDiscoveryServiceServer(grpcServer, xdsServer)
	v2.RegisterEndpointDiscoveryServiceServer(grpcServer, xdsServer)
	v2.RegisterClusterDiscoveryServiceServer(grpcServer, xdsServer)
	v2.RegisterRouteDiscoveryServiceServer(grpcServer, xdsServer)
	v2.RegisterListenerDiscoveryServiceServer(grpcServer, xdsServer)

	return
}

// HandleAbortSignal listens abort signals from OS and closes Fanplane gracefully
func (srv *Server) HandleAbortSignal(stop <-chan struct{}) (err error) {
	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("Shutting down gRPC server...")
			srv.grpcServer.GracefulStop()
		}
	}()

	return
}

//Start starts all components of the Fanplane discovery service on the port specified in FanplaneOptions
//Serving can be canceled at any time by closing the provided stop channel.
func (srv *Server) Start(stop <-chan struct{}) (err error) {
	// Now start all of the components.
	for _, fn := range srv.startFuncs {
		if err := fn(stop); err != nil {
			return err
		}
	}

	log.WithFields(log.Fields{"port": srv.Config.Port}).Info("management server listening")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", srv.Config.Port))
	if err != nil {
		log.WithError(err).Fatal("failed to listen")
	}

	go func() {
		<-stop
		srv.grpcServer.GracefulStop()
		log.Debugf("Fanplane server terminated: %v", err)
	}()

	if err = srv.grpcServer.Serve(lis); err != nil {
		log.WithError(err).Fatal(err)
	}

	return
}
