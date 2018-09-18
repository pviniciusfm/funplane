package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "fanplane",
		Short: "Fanplane agent.",
		Long:  "Fanplane agent runs is a simple way to provide centralized configuration for envoy sidecars .",
		RunE: func(c *cobra.Command, args []string) (err error) {
			fmt.Println("First command!")
			return
		},
	}
)

// func parseApplicationPorts() ([]uint16, error) {
// 	parsedPorts := make([]uint16, len(applicationPorts))
// 	for _, port := range applicationPorts {
// 		port := strings.TrimSpace(port)
// 		if len(port) > 0 {
// 			parsedPort, err := strconv.ParseUint(port, 10, 16)
// 			if err != nil {
// 				return nil, err
// 			}
// 			parsedPorts = append(parsedPorts, uint16(parsedPort))
// 		}
// 	}
// 	return parsedPorts, nil
// }

// func timeDuration(dur *types.Duration) time.Duration {
// 	out, err := types.DurationFromProto(dur)
// 	if err != nil {
// 		log.Warn(err)
// 	}
// 	return out
// }

func init() {
	// proxyCmd.PersistentFlags().StringVar((*string)(&registry), "serviceregistry",
	// 	string(serviceregistry.KubernetesRegistry),
	// 	fmt.Sprintf("Select the platform for service registry, options are {%s, %s, %s, %s, %s}",
	// 		serviceregistry.KubernetesRegistry, serviceregistry.ConsulRegistry,
	// 		serviceregistry.CloudFoundryRegistry, serviceregistry.MockRegistry, serviceregistry.ConfigRegistry))
	// proxyCmd.PersistentFlags().StringVar(&role.IPAddress, "ip", "",
	// 	"Proxy IP address. If not provided uses ${INSTANCE_IP} environment variable.")
	// proxyCmd.PersistentFlags().StringVar(&role.ID, "id", "",
	// 	"Proxy unique ID. If not provided uses ${POD_NAME}.${POD_NAMESPACE} from environment variables")
	// proxyCmd.PersistentFlags().StringVar(&role.Domain, "domain", "",
	// 	"DNS domain suffix. If not provided uses ${POD_NAMESPACE}.svc.cluster.local")
	// proxyCmd.PersistentFlags().Uint16Var(&statusPort, "statusPort", 0,
	// 	"HTTP Port on which to serve pilot agent status. If zero, agent status will not be provided.")
	// proxyCmd.PersistentFlags().StringSliceVar(&applicationPorts, "applicationPorts", []string{},
	// 	"Ports exposed by the application. Used to determine that Envoy is configured and ready to receive traffic.")

	// // Flags for proxy configuration
	// values := model.DefaultProxyConfig()
	// proxyCmd.PersistentFlags().StringVar(&configPath, "configPath", values.ConfigPath,
	// 	"Path to the generated configuration file directory")
	// proxyCmd.PersistentFlags().StringVar(&binaryPath, "binaryPath", values.BinaryPath,
	// 	"Path to the proxy binary")
	// proxyCmd.PersistentFlags().StringVar(&serviceCluster, "serviceCluster", values.ServiceCluster,
	// 	"Service cluster")
	// proxyCmd.PersistentFlags().StringVar(&availabilityZone, "availabilityZone", values.AvailabilityZone,
	// 	"Availability zone")
	// proxyCmd.PersistentFlags().DurationVar(&drainDuration, "drainDuration",
	// 	timeDuration(values.DrainDuration),
	// 	"The time in seconds that Envoy will drain connections during a hot restart")
	// proxyCmd.PersistentFlags().DurationVar(&parentShutdownDuration, "parentShutdownDuration",
	// 	timeDuration(values.ParentShutdownDuration),
	// 	"The time in seconds that Envoy will wait before shutting down the parent process during a hot restart")
	// proxyCmd.PersistentFlags().StringVar(&discoveryAddress, "discoveryAddress", values.DiscoveryAddress,
	// 	"Address of the discovery service exposing xDS (e.g. istio-pilot:8080)")
	// proxyCmd.PersistentFlags().DurationVar(&discoveryRefreshDelay, "discoveryRefreshDelay",
	// 	timeDuration(values.DiscoveryRefreshDelay),
	// 	"Polling interval for service discovery (used by EDS, CDS, LDS, but not RDS)")
	// proxyCmd.PersistentFlags().StringVar(&zipkinAddress, "zipkinAddress", values.ZipkinAddress,
	// 	"Address of the Zipkin service (e.g. zipkin:9411)")
	// proxyCmd.PersistentFlags().DurationVar(&connectTimeout, "connectTimeout",
	// 	timeDuration(values.ConnectTimeout),
	// 	"Connection timeout used by Envoy for supporting services")
	// proxyCmd.PersistentFlags().StringVar(&statsdUDPAddress, "statsdUdpAddress", values.StatsdUdpAddress,
	// 	"IP Address and Port of a statsd UDP listener (e.g. 10.75.241.127:9125)")
	// proxyCmd.PersistentFlags().Uint16Var(&proxyAdminPort, "proxyAdminPort", uint16(values.ProxyAdminPort),
	// 	"Port on which Envoy should listen for administrative commands")
	// proxyCmd.PersistentFlags().StringVar(&controlPlaneAuthPolicy, "controlPlaneAuthPolicy",
	// 	values.ControlPlaneAuthPolicy.String(), "Control Plane Authentication Policy")
	// proxyCmd.PersistentFlags().StringVar(&customConfigFile, "customConfigFile", values.CustomConfigFile,
	// 	"Path to the custom configuration file")
	// // Log levels are provided by the library https://github.com/gabime/spdlog, used by Envoy.
	// proxyCmd.PersistentFlags().StringVar(&proxyLogLevel, "proxyLogLevel", "warn",
	// 	fmt.Sprintf("The log level used to start the Envoy proxy (choose from {%s, %s, %s, %s, %s, %s, %s})",
	// 		"trace", "debug", "info", "warn", "err", "critical", "off"))
	// proxyCmd.PersistentFlags().IntVar(&concurrency, "concurrency", int(values.Concurrency),
	// 	"number of worker threads to run")
	// proxyCmd.PersistentFlags().BoolVar(&bootstrapv2, "bootstrapv2", true,
	// 	"Use bootstrap v2 - DEPRECATED")
	// proxyCmd.PersistentFlags().StringVar(&templateFile, "templateFile", "",
	// 	"Go template bootstrap config")
	// proxyCmd.PersistentFlags().BoolVar(&disableInternalTelemetry, "disableInternalTelemetry", false,
	// 	"Disable internal telemetry")

	// // Attach the Istio logging options to the command.
	// loggingOptions.AttachCobraFlags(rootCmd)

	// cmd.AddFlags(rootCmd)

	// rootCmd.AddCommand(proxyCmd)
	// rootCmd.AddCommand(version.CobraCommand())

	// rootCmd.AddCommand(collateral.CobraCommand(rootCmd, &doc.GenManHeader{
	// 	Title:   "Istio Pilot Agent",
	// 	Section: "pilot-agent CLI",
	// 	Manual:  "Istio Pilot Agent",
	// }))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(-1)
	}
}
