package adapter

import (
	"fmt"
	"github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.frg.tech/cloud/fanplane/pkg/apis/fanplane/v1alpha1"
	log "github.com/sirupsen/logrus"
)

type envoyBootstrapAdapter struct {
	bootstrap *v1alpha1.EnvoyBootstrap
}

func (adapter *envoyBootstrapAdapter) GetFanplaneObject() v1alpha1.FanplaneObject {
	return adapter.bootstrap
}

func (adapter *envoyBootstrapAdapter) BuildListeners() (result []cache.Resource, err error) {
	bootstrap := adapter.bootstrap.Spec.(*v2.Bootstrap)
	if err != nil {
		err = fmt.Errorf("couldn't process envoybootstrap model %s", err )
		return
	}

	result = []cache.Resource{}
	for _, listeners := range bootstrap.StaticResources.Listeners {
		result = append(result, &listeners)
	}

	return
}

func (envoyBootstrapAdapter) BuildRoutes() (result []cache.Resource, err error) {
	return
}

//TODO: Not efficient design as it is unwrapping it two times on BuildListeners and in BuildClusters
func (adapter *envoyBootstrapAdapter) BuildClusters() (result []cache.Resource, err error) {
	bootstrap := adapter.bootstrap.Spec.(*v2.Bootstrap)

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

func NewEnvoyBootstrapAdapter(obj v1alpha1.FanplaneObject) Adapter {
	result := &envoyBootstrapAdapter{}
	result.bootstrap = obj.(*v1alpha1.EnvoyBootstrap)
	if _, ok := obj.GetSpec().(*v2.Bootstrap); !ok {
		parsedBootstrap, err := v1alpha1.ParseEnvoyConfig(result.bootstrap.Spec)
		if err != nil {
			log.Fatal("couldn't parse EnvoyBootstrap in adapter", err)
		}
		result.bootstrap.Spec = parsedBootstrap
	}

	return result
}
