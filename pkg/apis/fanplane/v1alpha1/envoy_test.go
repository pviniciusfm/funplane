package v1alpha1

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
	"github.frg.tech/cloud/fanplane/pkg/apis/fanplane/test"
)

func TestEnvoyRawConversion(t *testing.T) {

	var fixtureTable = []test.FixtureTable{
		{
			"ValidEnvoyConfig",
			"testdata/envoy/envoy.yaml",
			"",
			test.Success,
		},
		{
			"ValidInvalidEnvoyStructure",
			"testdata/envoy/invalid-structure.yaml",
			"unknown value \"ROUND_ROBIN2\" for enum envoy.api.v2.Cluster_LbPolicy",
			test.Fail,
		},
		{
			"ValidInvalidEnvoyField",
			"testdata/envoy/invalid-fields.yaml",
			"unknown field \"socketaddress\"",
			test.Fail,
		},
	}

	for _, tt := range fixtureTable {
		testCaseName := fmt.Sprintf("Test%s", tt.Title)
		t.Run(testCaseName, func(t *testing.T) {
			checkEnvoyModel(t, tt)
		})
	}

}

func checkEnvoyModel(t *testing.T, table test.FixtureTable) {
	// Given: Reads valid envoy yaml file
	crd := &EnvoyBootstrap{}
	in, err := ioutil.ReadFile(table.InputFile)

	if assert.Nil(t, err) {
		return
	}
	// When: Unmarshal yaml
	// Unmarshal to CRD Type

	err = yaml.Unmarshal(in, crd)
	assert.Nil(t, err)
	assert.NotEmpty(t, crd)

	// Then: Should return a valid envoyBtstrp entity
	envoyBtstrp, err := ParseEnvoyConfig(crd.Spec)
	if table.ConversionResult == test.Fail {
		assert.Error(t, err, table.ExpectedOutputMessage)
		return
	}
	assert.Nil(t, err)
	assert.NotEmpty(t, envoyBtstrp)
	assert.Equal(t, "listener_0", envoyBtstrp.StaticResources.Listeners[0].Name)
	assert.Equal(t, "0.0.0.0", envoyBtstrp.StaticResources.Listeners[0].Address.GetSocketAddress().Address)
	assert.Equal(t, uint32(8800), envoyBtstrp.StaticResources.Listeners[0].Address.GetSocketAddress().GetPortValue())
}
