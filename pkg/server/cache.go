package server

import (
	"fmt"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyCache "github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/hashicorp/go-multierror"
	log "github.com/sirupsen/logrus"
	"github.frg.tech/cloud/fanplane/pkg/adapter"
	"github.frg.tech/cloud/fanplane/pkg/apis/fanplane/v1alpha1"
	"go.uber.org/atomic"
)

type Cache struct {
	envoyCache.SnapshotCache
	version atomic.Int32
}

func NewCache() Cache {
	result := Cache{}
	result.SnapshotCache = envoyCache.NewSnapshotCache(true, NodeHash{}, log.StandardLogger())
	return result
}

// NodeHash returns node ID as an ID
type NodeHash struct{}

// ID function
func (h NodeHash) ID(node *core.Node) string {
	if node == nil {
		return "unknown"
	}
	return node.Id
}

func (c *Cache) Add(obj v1alpha1.FanplaneObject) (err error) {
	cacheAdapter := adapter.NewAdapter(obj)
	rds, err := cacheAdapter.BuildRoutes()
	if err != nil {
		log.WithError(err).Errorln("couldn't parse routes.")
		multierror.Append(err, err)
	}

	cds, err := cacheAdapter.BuildClusters()
	if err != nil {
		log.WithError(err).Errorln("couldn't parse clusters.")
		return err
	}

	lds, err := cacheAdapter.BuildListeners()
	if err != nil {
		log.WithError(err).Errorln("couldn't parse listeners.")
		return err
	}

	snap := envoyCache.NewSnapshot(fmt.Sprint(c.version.Inc()), nil, cds, rds, lds)
	err = c.SnapshotCache.SetSnapshot(obj.GetObjectMeta().Name, snap)
	if err != nil {
		log.WithError(err).Errorln("couldn't set new config version.")
	}

	log.Infof("new version added in snapshot successfully.")

	return
}

func (c *Cache) Remove(obj v1alpha1.FanplaneObject) {
	c.SnapshotCache.ClearSnapshot(obj.GetObjectMeta().Name)
}

func (c *Cache) RemoveById(sideCarId string) {
	c.SnapshotCache.ClearSnapshot(sideCarId)
}