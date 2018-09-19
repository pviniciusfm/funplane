package v1alpha1

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.frg.tech/cloud/fanplane/pkg/apis/fanplane/test"
	"testing"
)

func TestGatewayModel(t *testing.T) {

	var fixtureTable = []test.FixtureTable{
		{
			"ValidGatewayConfig",
			"testdata/gateway/gateway.yaml",
			"",
			test.Success,
		},
		{
			"InvalidGatewayStructure",
			"testdata/gateway/invalid-structure.yaml",
			"invalid config: Spec.Listener: zero value",
			test.Fail,
		},
		{
			"InvalidGatewayField",
			"testdata/gateway/invalid-fields.yaml",
			"parsing config: error unmarshaling JSON: " +
				"invalid RouteType \"CONSUL_ERROR\"",
			test.Fail,
		},
	}

	for _, table := range fixtureTable {
		testCaseName := fmt.Sprintf("Test%s", table.Title)
		t.Run(testCaseName, func(t *testing.T) {
			checkGatewayModel(t, table)
		})
	}

}

func checkGatewayModel(t *testing.T, table test.FixtureTable) {
	// Given: An gateway configuration file
	// When: LoadGateway
	crd, err := LoadGateway(table.InputFile)

	if table.ConversionResult == test.Fail {
		//If invalid conversion should have expected output message
		assert.EqualError(t, err, table.ExpectedOutputMessage)
		return
	}

	// Then: Verify valid content
	if assert.NotEmpty(t, crd) {
		assert.Equal(t, crd.Spec.Listener.Protocol, HTTPS)
		assert.NotEmpty(t, crd.Spec.Routes)
		assert.NotZero(t, crd.Spec.Routes[0].RouteType)
		assert.NotEmpty(t, crd.Spec.Routes[0].TrafficPolicy)
		assert.NotZero(t, crd.Spec.Routes[0].TrafficPolicy.LoadBalancerSettings.LoadBalancerType)
		assert.NotEmpty(t, crd.Spec.Routes[0].TrafficPolicy.RetryPolicy)
		assert.NotEmpty(t, crd.Spec.Routes[0].TrafficPolicy.ConnectionPoolSettings)
		assert.NotEmpty(t, crd.Spec.Selector)
	}
}
