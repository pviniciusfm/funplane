package adapter

import (
	"fmt"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.frg.tech/cloud/fanplane/pkg/apis/fanplane/v1alpha1"
)

type EnvoyBootstrapAdapter struct {
	Aggregated
}

func (EnvoyBootstrapAdapter) GetListeners(object v1alpha1.FanplaneObject) (result []cache.Resource, err error) {
	bootstrap, err := v1alpha1.ParseEnvoyConfig(object.GetSpec())
	if err != nil {
		err = fmt.Errorf("couldn't process envoybootstrap model %s", err )
		return
	}

	result = []cache.Resource{}
	for _, listeners := range bootstrap.GetStaticResources().GetListeners() {
		result = append(result, &listeners)
	}

	return
}

func (EnvoyBootstrapAdapter) GetRoutes(object v1alpha1.FanplaneObject) (result []cache.Resource, err error) {
	return
}

func (EnvoyBootstrapAdapter) GetEndpoints(object v1alpha1.FanplaneObject) (result []cache.Resource, err error) {
	return
}

//TODO: Not efficient design is unwrapping it two times on GetListeners and in GetClusters
func (EnvoyBootstrapAdapter) GetClusters(object v1alpha1.FanplaneObject) (result []cache.Resource, err error) {
	bootstrap, err := v1alpha1.ParseEnvoyConfig(object.GetSpec())

	if err != nil {
		err = fmt.Errorf("couldn't process envoybootstrap model %s", err )
		return
	}

	result = []cache.Resource{}
	for _, clusters := range bootstrap.GetStaticResources().GetClusters() {
		result = append(result, &clusters)
	}

	return
}
