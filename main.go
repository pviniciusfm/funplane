package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.frg.tech/cloud/fanplane/pkg/server"
	"os"
	"time"
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

func init() {
	defaultDrainDuration, _ := time.ParseDuration("1s")
	defaultDiscoveryRefreshDelay, _ := time.ParseDuration("1s")

	config := &server.FanplaneConfig{}
	rootCmd.PersistentFlags().IntVar(&config.Port, "port", 1800, "Fanplane listen port")
	rootCmd.PersistentFlags().StringVar(&config.IPAddress, "ip", "",
		"Proxy IP address. If not provided uses ${INSTANCE_IP} environment variable.")
	rootCmd.PersistentFlags().StringVar(&config.ID, "id", "",
		"Proxy unique ID. If not provided uses ${POD_NAME}.${POD_NAMESPACE} from environment variables")
	rootCmd.PersistentFlags().StringVar(&config.Domain, "domain", "",
		"DNS domain suffix. Used for consul service discovery")

	rootCmd.PersistentFlags().StringVar(&config.ConfigPath, "configPath", config.ConfigPath,
		"Path to the generated configuration file directory")
	rootCmd.PersistentFlags().DurationVar(&config.DrainDuration, "drainDuration", defaultDrainDuration,
		"The time in seconds that Envoy will drain connections during a hot restart")

	rootCmd.PersistentFlags().DurationVar(&config.DiscoveryRefreshDelay, "discoveryRefreshDelay",
		defaultDiscoveryRefreshDelay,
		"Polling interval for service discovery (used by EDS, CDS, LDS, but not RDS)")

	rootCmd.PersistentFlags().StringVar(&config.StatsdUDPAddress, "statsdUdpAddress", "localhost:8125",
		"IP Address and Port of a Statsd UDP listener (e.g. 10.75.241.127:9125)")

	rootCmd.PersistentFlags().StringVar(&config.LogLevel, "proxyLogLevel", "warn",
		fmt.Sprintf("The log level used to start the Envoy proxy (choose from {%s, %s, %s, %s, %s, %s, %s})",
			"trace", "debug", "info", "warn", "err", "critical", "off"))
	rootCmd = &cobra.Command{
		Use:   "fanplane",
		Short: "Fanplane agent.",
		Long:  "Fanplane agent runs is a simple way to provide centralized configuration for envoy sidecars .",
		RunE: func(c *cobra.Command, args []string) (err error) {
			signal := make(chan struct{})

			ctx := context.Background()
			server.NewManagementServer(ctx, config)

			<-signal

			return
		},
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(-1)
	}
}
