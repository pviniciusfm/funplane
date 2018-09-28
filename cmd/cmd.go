package cmd

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.frg.tech/cloud/fanplane/pkg/registry/kube"
	"github.frg.tech/cloud/fanplane/pkg/server"
	"time"
)

var (
	FanplaneCmd *cobra.Command         = &cobra.Command{}
	config      *server.FanplaneConfig = &server.FanplaneConfig{
		Port:     18000,
		ADSMode:  true,
		Domain:   "lab.frgcloud.com",
		ID:       "fanplane",
		LogLevel: "debug",
	}
	gitSha1 string
	version string
)

func init() {
	viper.SetConfigName("fanplane")
	FanplaneCmd.AddCommand(docCmd)

	FanplaneCmd.Version = fmt.Sprintf("%s {GitSha1: %s}", version, gitSha1)

	defaultDrainDuration, _ := time.ParseDuration("1s")
	defaultDiscoveryRefreshDelay, _ := time.ParseDuration("1s")

	config := &server.FanplaneConfig{}
	FanplaneCmd.PersistentFlags().StringVar(&config.KubeCfgFile, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	FanplaneCmd.PersistentFlags().StringVar(&config.MasterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	FanplaneCmd.PersistentFlags().IntVar(&config.Port, "port", 18000, "Fanplane listen port")
	FanplaneCmd.PersistentFlags().StringVar(&config.IPAddress, "ip", "0.0.0.0",
		"Proxy IP address. If not provided uses ${INSTANCE_IP} environment variable.")
	FanplaneCmd.PersistentFlags().StringVar(&config.ID, "id", "",
		"Proxy unique ID. If not provided uses ${POD_NAME}.${POD_NAMESPACE} from environment variables")
	FanplaneCmd.PersistentFlags().StringVar(&config.Domain, "domain", "",
		"DNS domain suffix. Used for consul service discovery")

	FanplaneCmd.PersistentFlags().StringVar(&config.ConfigPath, "configPath", config.ConfigPath,
		"Path to the generated configuration file directory")
	FanplaneCmd.PersistentFlags().DurationVar(&config.DrainDuration, "drainDuration", defaultDrainDuration,
		"The time in seconds that Envoy will drain connections during a hot restart")

	FanplaneCmd.PersistentFlags().DurationVar(&config.DiscoveryRefreshDelay, "discoveryRefreshDelay",
		defaultDiscoveryRefreshDelay,
		"Polling interval for service discovery (used by EDS, CDS, LDS, but not RDS)")

	FanplaneCmd.PersistentFlags().StringVar(&config.StatsdUDPAddress, "statsdUdpAddress", "localhost:8125",
		"IP Address and Port of a Statsd UDP listener (e.g. 10.75.241.127:9125)")

	FanplaneCmd.PersistentFlags().StringVar(&config.LogLevel, "proxyLogLevel", "warn",
		fmt.Sprintf("The log level used to start the Envoy proxy (choose from {%s, %s, %s, %s, %s, %s, %s})",
			"trace", "debug", "info", "warn", "err", "critical", "off"))
	FanplaneCmd = &cobra.Command{
		Use:   "fanplane",
		Short: "Fanplane agent.",
		Long:  "Fanplane agent runs is a simple way to provide centralized configuration for envoy sidecars .",
		RunE: func(c *cobra.Command, args []string) (err error) {
			signal := make(chan struct{})

			ctx := context.Background()
			sv, err := server.NewManagementServer(ctx, config)
			go kube.Initialize(config, sv.Cache)

			if err != nil {
				log.WithError(err).Fatal("Server stopped abruptly")
				return
			}

			err = sv.Start(signal)
			if err != nil {
				log.WithError(err).Fatal("Couldn't initialize server")
			}

			return
		},
	}

}
