package adapter

import (
	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/spf13/viper"
	testAssert "github.com/stretchr/testify/assert"
	"github.frg.tech/cloud/fanplane/pkg/apis/fanplane/v1alpha1"
	"testing"
)

func init() {
	viper.Set("domain", "dev.frgcloud.com")
}

func TestConvertGateway(t *testing.T) {
	assert := testAssert.New(t)
	//Given: A valid gateway entity
	crd, err := v1alpha1.LoadGateway("testdata/gateway-conversion.yaml")
	if !assert.Nil(err) {
		return
	}

	adapter := NewAdapter(crd)

	//When: Convert listeners
	listeners, err := adapter.BuildListeners()
	if err != nil {
		t.Fatalf("error converting gateway %s", err)
	}

	//Then: Should output valid Listeners and Clusters entities of envoy's domain
	assert.NotNil(listeners)

	listener := listeners[0].(*v2.Listener)
	assert.NotZero(listener.Name)
	assert.NotZero(listener.Address.GetSocketAddress().Address)
	assert.NotZero(listener.GetFilterChains()[0].Filters[0].Config)

	clusters, err := adapter.BuildClusters()
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

	//Given: A valid gateway entity
	crd, err := v1alpha1.LoadEnvoyBootstrap("testdata/bootstrap.yaml")
	if !assert.Nil(err) {
		return
	}

	adapter := NewAdapter(crd)

	//When: Convert listeners
	listeners, err := adapter.BuildListeners()
	if err != nil {
		t.Fatalf("error converting gateway %s", err)
	}

	//Then: Should output valid envoy LDS response
	assert.NotNil(listeners)

	listener := listeners[0].(*v2.Listener)
	assert.NotZero(listener.Name)
	assert.NotZero(listener.Address.GetSocketAddress().Address)
	assert.NotZero(listener.GetFilterChains()[0].Filters[0].Config)

	//When: Convert clusters
	clusters, err := adapter.BuildClusters()
	if err != nil {
		t.Fatalf("error converting gateway %s", err)
	}

	//Then: Should output valid envoy CDS response
	assert.NotEmpty(clusters)
	cluster := clusters[0].(*v2.Cluster)
	assert.NotZero(cluster.Name)
	assert.NotZero(cluster.TlsContext.Sni)
}
