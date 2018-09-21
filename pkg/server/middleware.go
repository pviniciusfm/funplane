package server

import (
	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Callback struct {
	signal   chan struct{}
	fetches  int
	requests int
	mu       sync.Mutex
}

// OnStreamOpen is called once an xDS stream is open with a stream ID and the type URL (or "" for ADS).
func (*Callback) OnStreamOpen(id int64, typ string) {
	streamType := typ
	if typ == ""{
		streamType = "ads"
	}
	log.Debugf("stream opened %d open for %s", id, streamType)
}

// OnStreamClosed is called immediately prior to closing an xDS stream with a stream ID.
func (*Callback) OnStreamClosed(id int64) {
	log.Debugf("stream closed %d closed", id)
}

// OnStreamRequest is called once a request is received on a stream.
func (*Callback) OnStreamRequest(int64, *v2.DiscoveryRequest) {}


// OnStreamResponse is called immediately prior to sending a response on a stream.
func (cb *Callback) OnStreamResponse(id int64, req *v2.DiscoveryRequest, resp *v2.DiscoveryResponse) {
	log.Debugf("stream response %s", req)
}

// OnFetchRequest is called for each Fetch request
func (*Callback) OnFetchRequest(*v2.DiscoveryRequest) {}

// OnFetchResponse is called immediately prior to sending a response.
func (*Callback) OnFetchResponse(*v2.DiscoveryRequest, *v2.DiscoveryResponse) {}
