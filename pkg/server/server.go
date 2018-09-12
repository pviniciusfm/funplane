package server

// Server contains the runtime configuration for the fanplane control service.
type Server struct {
}

// NewServer creates a new Server instance based on the provided arguments.
// func NewServer(args PilotArgs) (*Server, error) {

// }

// Start starts all components of the Pilot discovery service on the port specified in DiscoveryServiceOptions.
// If Port == 0, a port number is automatically chosen. Content serving is started by this method,
// but is executed asynchronously. Serving can be cancelled at any time by closing the provided stop channel.
// func (s *Server) Start(stop <-chan struct{}) error {
// 	// Now start all of the components.
// 	for _, fn := range s.startFuncs {
// 		if err := fn(stop); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// // initKubeClient creates the k8s client if running in an k8s environment.
// func (s *Server) initKubeClient(args *PilotArgs) error {
// 	if hasKubeRegistry(args) && args.Config.FileDir == "" {
// 		client, kuberr := kubelib.CreateClientset(s.getKubeCfgFile(args), "")
// 		if kuberr != nil {
// 			return multierror.Prefix(kuberr, "failed to connect to Kubernetes API.")
// 		}
// 		s.kubeClient = client

// 	}

// 	return nil
// }
