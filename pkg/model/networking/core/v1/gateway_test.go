package model

import (
	"testing"
)

const (
	ValidGatewayPath = "testdata/gateway.yaml"
)

func TestGatewayMarshall(t *testing.T) {
	err := LoadGateway(ValidGatewayPath)
	assert.AssertNotNil(err)
}
