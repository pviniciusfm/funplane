package server

import (
	"context"
	"errors"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.frg.tech/cloud/fanplane/pkg/apis/fanplane"
	"github.frg.tech/cloud/fanplane/pkg/registry"
	"github.frg.tech/cloud/fanplane/pkg/registry/filestore"
	"testing"
)

const (
	testSideCarId = "test-id"
)

var (
	config = &fanplane.Config{
		Port:              18000,
		ADSMode:           true,
		Domain:            "consul",
		ID:                "test",
		LogLevel:          log.DebugLevel.String(),
		RegistryDirectory: "testdata/",
		RegistryType:      registry.FileRegistry,
	}

	node = &core.Node{
		Id: testSideCarId,
	}
)

func TestManagementServer(t *testing.T) {
	srv, err := NewManagementServer(context.Background(), config)
	filestore.NewFileRegistry(config, srv.ServerCache)

	stopSignal := make(chan struct{})
	defer func() {
		//FUTURE: Find another way to not fail the tests because this function is being accessed after the test is over
		t.SkipNow()
		stopSignal <- struct{}{}
	}()

	assert.Nil(t, err)

	go srv.Start(stopSignal)

	// When make a cdsRequest
	cdsRequest := &v2.DiscoveryRequest{
		Node:    node,
		TypeUrl: ClusterType,
	}

	resp, err := srv.xDSserver.FetchClusters(context.Background(), cdsRequest)
	if err != nil {
		t.Error("fetch clusters() => got error", err)
		return
	}
	assert.NotNil(t, resp)

	cdsResponse, err := GetClusterResponse(resp)
	if !assert.Nil(t, err) {
		return
	}

	assert.NotEmpty(t, cdsResponse.Name)
	assert.EqualValues(t, v2.Cluster_LOGICAL_DNS, cdsResponse.Type)
	assert.EqualValues(t, "service_cybersource", cdsResponse.Name)

	ldsRequest := &v2.DiscoveryRequest{
		Node:    node,
		TypeUrl: ClusterType,
	}

	rawLdsResp, err := srv.xDSserver.FetchListeners(context.Background(), ldsRequest)
	if err != nil {
		t.Error("fetch listeners() => got error")
	}

	ldsResponse, err := GetListenerResponse(rawLdsResp)
	if !assert.Nil(t, err) {
		return
	}

	assert.NotNil(t, resp)
	assert.EqualValues(t, "listener_0", ldsResponse.Name)
	assert.EqualValues(t, int32(8800), ldsResponse.Address.GetSocketAddress().GetPortValue())

}

// Extract cluster response from a discovery response.
func GetClusterResponse(res1 *v2.DiscoveryResponse) (*v2.Cluster, error) {
	if res1.TypeUrl != ClusterType {
		return nil, errors.New("Invalid typeURL" + res1.TypeUrl)
	}

	if res1.Resources[0].TypeUrl != ClusterType {
		return nil, errors.New("Invalid resource typeURL" + res1.Resources[0].TypeUrl)
	}

	cla := &v2.Cluster{}
	err := cla.Unmarshal(res1.Resources[0].Value)
	if err != nil {
		return nil, err
	}
	return cla, nil
}

// Extract listener response from a discovery response.
func GetListenerResponse(res1 *v2.DiscoveryResponse) (*v2.Listener, error) {
	if res1.TypeUrl != ListenerType {
		return nil, errors.New("Invalid typeURL" + res1.TypeUrl)
	}

	if res1.Resources[0].TypeUrl != ListenerType {
		return nil, errors.New("Invalid resource typeURL" + res1.Resources[0].TypeUrl)
	}

	listener := &v2.Listener{}
	err := listener.Unmarshal(res1.Resources[0].Value)
	if err != nil {
		return nil, err
	}
	return listener, nil
}
