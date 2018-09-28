//Adapter package is responsible for the Translation of Fanplane Objects to the XDS Envoy Models
package adapter

import (
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.frg.tech/cloud/fanplane/pkg/apis/fanplane/v1alpha1"
	log "github.com/sirupsen/logrus"
)

type Aggregated interface {
	// GetListeners returns envoy listeners messages from listener FanplaneObject. This is the LDS output
	GetListeners(object v1alpha1.FanplaneObject) ([]cache.Resource, error)
	// GetRoutes returns envoy Route messages from a FanplaneObject. This is the RDS output
	GetRoutes(object v1alpha1.FanplaneObject) ([]cache.Resource, error)
	// GetEndpoints returns envoy Endpoint messages from a FanplaneObject. This is the EDS output
	GetEndpoints(object v1alpha1.FanplaneObject) ([]cache.Resource, error)
	// GetClusters returns the list of clusters for the given FanplaneObject. This is the CDS output
	// For outbound: Cluster for each service/subset hostname or CIDR with SNI set to service hostname
	// Cluster type based on resolution
	GetClusters(object v1alpha1.FanplaneObject) ([]cache.Resource, error)
}

func AddCachedResources(adapter Aggregated,
	object v1alpha1.FanplaneObject, version string, snapshotCache cache.SnapshotCache) (err error) {

	rds, err := adapter.GetRoutes(object)
	if err != nil {
		log.WithError(err).Errorln("couldn't parse routes.")
		return err
	}

	cds, err := adapter.GetClusters(object)
	if err != nil {
		log.WithError(err).Errorln("couldn't parse clusters.")
		return err
	}

	eds, err := adapter.GetEndpoints(object)
	if err != nil {
		log.WithError(err).Errorln("couldn't parse endpoints.")
		return err
	}

	lds, err := adapter.GetListeners(object)
	if err != nil {
		log.WithError(err).Errorln("couldn't parse listeners.")
		return err
	}

	snap := cache.NewSnapshot(version, eds, cds, rds, lds)
	err = snapshotCache.SetSnapshot(object.GetSidecarSelector(), snap)
	if err != nil{
		log.WithError(err).Errorln("couldn't set new config version.")
	}

	log.Infof("new version added in snapshot successfully.")

	return
}
