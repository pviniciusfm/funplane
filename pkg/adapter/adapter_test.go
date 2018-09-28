package adapter

import (
	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/ghodss/yaml"
	"github.com/spf13/viper"
	testAssert "github.com/stretchr/testify/assert"
	"github.frg.tech/cloud/fanplane/pkg/apis/fanplane/v1alpha1"
	"io/ioutil"
	"testing"
)

func TestConvertGateway(t *testing.T) {
	assert := testAssert.New(t)
	//Given: A valid gateway entity
	viper.Set("domain", "dev.frgcloud.com")
	adapter := &GatewayAdapter{}
	crd := &v1alpha1.Gateway{}
	inputFile := "testdata/gateway-conversion.yaml"
	in, err := ioutil.ReadFile(inputFile)
	if !assert.Nil(err) {
		return
	}

	err = yaml.Unmarshal(in, crd)
	if !assert.Nil(err) {
		return
	}

	//When: Convert listeners
	listeners, err := adapter.GetListeners(crd)

	//Then: Should output valid Listeners and Clusters entities of envoy's domain
	assert.NotNil(listeners)

	listener := listeners[0].(*v2.Listener)
	assert.NotZero(listener.Name)
	assert.NotZero(listener.Address.GetSocketAddress().Address)
	assert.NotZero(listener.GetFilterChains()[0].Filters[0].Config)

	clusters, err := adapter.GetClusters(crd)
	if !assert.Nil(err) {
		return
	}
	assert.NotEmpty(clusters)

	cluster := clusters[0].(*v2.Cluster)
	assert.NotZero(cluster.Name)

	cluster = clusters[1].(*v2.Cluster)
	assert.NotZero(cluster.TlsContext.Sni)
}

func TestConvertEnvoyBootstrap(t *testing.T) {
	assert := testAssert.New(t)

	//Given: A Fanplane EnvoyBootstrap object
	adapter := &EnvoyBootstrapAdapter{}
	crd := &v1alpha1.EnvoyBootstrap{}
	inputFile := "testdata/bootstrap.yaml"

	in, err := ioutil.ReadFile(inputFile)
	if !assert.Nil(err) {
		return
	}

	err = yaml.Unmarshal(in, crd)
	if !assert.Nil(err) {
		return
	}

	//When: Convert listeners
	listeners, err := adapter.GetListeners(crd)

	//Then: Should output valid envoy LDS response
	assert.NotNil(listeners)

	listener := listeners[0].(*v2.Listener)
	assert.NotZero(listener.Name)
	assert.NotZero(listener.Address.GetSocketAddress().Address)
	assert.NotZero(listener.GetFilterChains()[0].Filters[0].Config)

	//When: Convert clusters
	clusters, err := adapter.GetClusters(crd)

	//Then: Should output valid envoy CDS response
	assert.NotEmpty(clusters)
	cluster := clusters[0].(*v2.Cluster)
	assert.NotZero(cluster.Name)
	assert.NotZero(cluster.TlsContext.Sni)
}
