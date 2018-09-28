package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/stretchr/testify/assert"
	"github.frg.tech/cloud/fanplane/pkg/apis/fanplane/v1alpha1"
	"io/ioutil"
	"testing"
)

const (
	testSideCarId   = "test-id"
	rawClusterYaml  = "testdata/raw-cluster.yaml"
	rawListenerYaml = "testdata/raw-listener.yaml"
)

var (
	config *FanplaneConfig = &FanplaneConfig{
		Port:     18000,
		ADSMode:  true,
		Domain:   "consul",
		ID:       "test",
		LogLevel: "debug",
	}

	node = &core.Node{
		Id: testSideCarId,
	}
)


func LoadMockRegistry(t *testing.T) (snap cache.Snapshot, error error) {
	rawYaml, err := ioutil.ReadFile(rawClusterYaml)
	if err != nil {
		return
	}

	googleCluster, err := v1alpha1.FromYAML(string(rawYaml), "envoy.api.v2.Cluster")
	if err != nil {
		t.Fatalf("couldn't set cluster entity in snapsoht. %s", err)
		return
	}

	rawYaml, err = ioutil.ReadFile(rawListenerYaml)
	testListener, err := v1alpha1.FromYAML(string(rawYaml), "envoy.api.v2.Listener")
	if err != nil {
		t.Fatalf("couldn't set listener entity in snapsoht. %s", err)
		return
	}

	cdsRes := googleCluster.(*v2.Cluster)
	ldsRes := testListener.(*v2.Listener)
	cds := []cache.Resource{cdsRes}
	lds := []cache.Resource{ldsRes}
	snap = cache.NewSnapshot(fmt.Sprint("1"), nil, cds, nil, lds)

	return
}

func TestManagementServer(t *testing.T) {
	snap, err := LoadMockRegistry(t)
	if err != nil {
		t.Fail()
	}
	stopSignal := make(chan struct{})
	defer func() { stopSignal <- struct{}{} }()

	srv, err := NewManagementServer(context.Background(), config)
	srv.Cache.SetSnapshot(testSideCarId, snap)
	assert.Nil(t, err)

	go srv.Start(stopSignal)

	// When make a cdsRequest
	cdsRequest := &v2.DiscoveryRequest{
		Node:    node,
		TypeUrl: ClusterType,
	}

	resp, err := srv.xDSserver.FetchClusters(context.Background(), cdsRequest)
	if err != nil {
		t.Error("fetch clusters() => got error")
	}
	assert.NotNil(t, resp)

	cdsResponse, err := GetClusterResponse(resp)
	if !assert.Nil(t, err) {
		return
	}

	assert.NotEmpty(t, cdsResponse.Name)
	assert.EqualValues(t, v2.Cluster_LEAST_REQUEST, cdsResponse.LbPolicy)
	assert.EqualValues(t, "google-test", cdsResponse.Name)

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
	assert.EqualValues(t,"listener_0", ldsResponse.Name)
	assert.EqualValues(t,int32(8800), ldsResponse.Address.GetSocketAddress().GetPortValue())
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
