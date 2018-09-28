package server

import "time"

//FanplaneConfig is the object that holds all parameterized configuration of envoy control plane
type FanplaneConfig struct {
	Port                  int
	IPAddress             string
	ID                    string
	Domain                string
	AgentStatusPort       int
	MasterURL             string
	ConfigPath            string
	DrainDuration         time.Duration
	DiscoveryRefreshDelay time.Duration
	StatsdUDPAddress      string
	LogLevel              string
	KubeCfgFile           string
	ADSMode               bool
}

//StatsDConfig holds the configuration required to contact StatsD server
type StatsDConfig struct {
	appName              string
	flushPeriodInSeconds int
	port                 int
	enabled              bool
}

func NewStatsDConfig(appName string, flushPeriodInSeconds, port int, enabled bool) *StatsDConfig {
	return &StatsDConfig{
		appName:              appName,
		flushPeriodInSeconds: flushPeriodInSeconds,
		port:                 port,
		enabled:              enabled,
	}
}

func (sdc *StatsDConfig) AppName() string {
	return sdc.appName
}

func (sdc *StatsDConfig) FlushPeriodInSeconds() int {
	return sdc.flushPeriodInSeconds
}

func (sdc *StatsDConfig) Port() int {
	return sdc.port
}

func (sdc *StatsDConfig) Enabled() bool {
	return sdc.enabled
}
