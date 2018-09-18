package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	ValidGatewayPath = "testdata/gateway.yaml"
)

func TestGatewayMarshall(t *testing.T) {
	gateway, err := LoadGateway(ValidGatewayPath)
	assert.NotNil(t, gateway)
	assert.Nil(t, err)
}
