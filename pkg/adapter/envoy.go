package adapter

import (
	"fmt"
	"github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	log "github.com/sirupsen/logrus"
	"github.frg.tech/cloud/fanplane/pkg/apis/fanplane"
	"github.frg.tech/cloud/fanplane/pkg/apis/fanplane/v1alpha1"
)

type envoyBootstrapAdapter struct {
	bootstrap *v1alpha1.EnvoyBootstrap
}

func (adapter *envoyBootstrapAdapter) GetFanplaneObject() fanplane.Object {
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

//NewEnvoyBootstrapAdapter returns an adapter from FanplaneObject to EnvoyBootstrap proto message
func NewEnvoyBootstrapAdapter(obj fanplane.Object) Adapter {
	result := &envoyBootstrapAdapter{}

	if _, ok := obj.(*fanplane.Kind); ok {
		result.bootstrap = &v1alpha1.EnvoyBootstrap{TypeMeta: obj.GetTypeMeta(), ObjectMeta: obj.GetObjectMeta()}
	}else{
		result.bootstrap = obj.(*v1alpha1.EnvoyBootstrap)
	}

	//Only converts object if not adapted
	if _, ok := obj.GetSpec().(*v2.Bootstrap); !ok {
		parsedBootstrap, err := v1alpha1.ParseEnvoyConfig(obj.GetSpec())
		if err != nil {
			log.Fatal("couldn't parse EnvoyBootstrap in adapter", err)
		}

		if _, ok := obj.(*v1alpha1.EnvoyBootstrap); ok{
			result.bootstrap = obj.(*v1alpha1.EnvoyBootstrap)
		} else {
			//Adapt from FanplaneKind
			result.bootstrap = &v1alpha1.EnvoyBootstrap{
				TypeMeta: obj.GetTypeMeta(),
				ObjectMeta: obj.GetObjectMeta(),
			}
		}

		//Adapt from EnvoyBootstrap
		result.bootstrap.Spec = parsedBootstrap
	}

	return result
}
