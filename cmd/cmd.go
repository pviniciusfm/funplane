package cmd

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.frg.tech/cloud/fanplane/pkg/apis/fanplane"
	"github.frg.tech/cloud/fanplane/pkg/cache"
	"github.frg.tech/cloud/fanplane/pkg/registry"
	"github.frg.tech/cloud/fanplane/pkg/registry/filestore"
	"github.frg.tech/cloud/fanplane/pkg/registry/kube"
	"github.frg.tech/cloud/fanplane/pkg/server"
	"os"
	"runtime/debug"
)

var (
	FanplaneCmd *cobra.Command = &cobra.Command{
		Use:   "fanplane",
		Short: "Fanplane agent.",
		Long: `
 _______           ______  _                    
(_______)         (_____ \| |                   
 _____ _____ ____  _____) ) | _____ ____  _____ 
|  ___|____ |  _ \|  ____/| |(____ |  _ \| ___ |
| |   / ___ | | | | |     | |/ ___ | | | | ____|
|_|   \_____|_| |_|_|      \_)_____|_| |_|_____)

Fanplane is a simple way to provide centralized configuration for envoy sidecars .

Use --logLevel trace to increase log verbosity   

Example of utilization:

fanplane --domain dynamic.dev.frgcloud.com --port 18000 --logLevel trace
`,
		RunE: func(c *cobra.Command, args []string) (err error) {
			initializeLog(config)
			initializeServer(config)
			return
		},
	}
	//Those are the config defaults
	config *fanplane.Config = &fanplane.Config{
		Port:              18000,
		IPAddress:         "0.0.0.0",
		Domain:            "consul",
		StatsdUDPAddress:  "127.0.0.1:8125",
		LogLevel:          "debug",
		KubeCfgFile:       "",
		ADSMode:           true,
		RegistryType:      "kube",
		RegistryDirectory: "registry/",
	}

	gitSha1 string
	version string
)

func init() {
	viper.SetConfigName("fanplane")
	FanplaneCmd.AddCommand(docCmd)

	FanplaneCmd.Version = fmt.Sprintf("%s {GitSha1: %s}", version, gitSha1)

	FanplaneCmd.Flags().StringVar(&config.KubeCfgFile, "kubeconfig", config.KubeCfgFile, "Path to a kubeconfig. Only required if out-of-cluster.")
	FanplaneCmd.Flags().StringVar(&config.MasterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	FanplaneCmd.Flags().IntVar(&config.Port, "port", config.Port, "Fanplane listen port")
	FanplaneCmd.Flags().StringVar(&config.IPAddress, "ip", config.IPAddress,
		"Proxy IP address. If not provided uses ${INSTANCE_IP} environment variable.")

	FanplaneCmd.Flags().StringVar(&config.Domain, "domain", config.Domain,
		"DNS domain suffix. Used for consul service discovery")

	FanplaneCmd.Flags().StringVar(&config.RegistryType, "registryType", config.RegistryType,
		fmt.Sprintf("which configuration registry will be used for the agent. Valid options %s", registry.AllRegistryTypes))

	FanplaneCmd.Flags().StringVar(&config.RegistryDirectory, "fileRegistryPath", config.RegistryDirectory,
		"Path to the generated configuration filestore directory")

	FanplaneCmd.Flags().StringVar(&config.StatsdUDPAddress, "statsdUdpAddress", config.StatsdUDPAddress,
		"IP Address and Port of a Statsd UDP listener (e.g. 10.75.241.127:9125)")

	FanplaneCmd.Flags().StringVar(&config.LogLevel, "logLevel", config.LogLevel,
		fmt.Sprintf("The log level used to start the Fanplane (choose from {%s, %s, %s, %s, %s, %s, %s})",
			"trace", "debug", "info", "warn", "err", "critical", "off"))

	viper.BindPFlag("domain", FanplaneCmd.Flags().Lookup("domain"))
}

//initializeRegistry starts config registry according with Fanplane configuration
func initializeRegistry(config *fanplane.Config, cache cache.Cache) error {
	log.Debug("Initializing registry")

	registryType, err := registry.ParseRegistryType(config.RegistryType)
	if err != nil {
		log.WithError(err).Fatal("couldn't initialize registry.")
	}

	switch registryType {
	case registry.FileRegistry:
		registry, err := filestore.NewFileRegistry(config, cache)
		if err != nil {
			return err
		}
		go registry.StartWatch()
		return nil
	case registry.KubeRegistry:
		go kube.Initialize(config, cache)
		return nil
	}
	return fmt.Errorf("registryType %q doesn not exist", registryType)
}

//initializeLog configures log level of the entire application
func initializeLog(config *fanplane.Config) {
	result, err := log.ParseLevel(config.LogLevel)

	if err != nil {
		log.Fatalf("couldn't find log level %s", config.LogLevel)
	}

	log.SetLevel(result)
}

//initializeServer initializes server and registry
func initializeServer(config *fanplane.Config) {
	signal := make(chan struct{})
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("fanplane internal error. %q", r)
			log.Debug(string(debug.Stack()))
			os.Exit(1)
		}
	}()

	ctx := context.Background()
	sv, err := server.NewManagementServer(ctx, config)
	if err != nil {
		log.WithError(err).Fatal("server stopped abruptly")
	}

	err = initializeRegistry(config, sv.ServerCache)
	if err != nil {
		log.WithError(err).Fatal("couldn't initialize registry", err)
	}

	sv.Start(signal)
	if err != nil {
		log.WithError(err).Fatal("server stopped abruptly")
	}
}
