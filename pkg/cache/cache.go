package cache

import (
	"fmt"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyCache "github.com/envoyproxy/go-control-plane/pkg/cache"
	log "github.com/sirupsen/logrus"
	"github.frg.tech/cloud/fanplane/pkg/adapter"
	"github.frg.tech/cloud/fanplane/pkg/apis/fanplane"
	"go.uber.org/atomic"
)

// NodeHash makes sidecar envoy ID the cache key
type NodeHash struct{}

// ID function uses node id as the unique id
func (h NodeHash) ID(node *core.Node) string {
	if node == nil {
		return "unknown"
	}
	return node.Id
}

//Cache wraps go-control-plane logic to integrate with Fanplane Object model
type Cache interface {
	GetVersion() *atomic.Int32
	Add(obj fanplane.Object) (err error)
	Remove(obj fanplane.Object)
	RemoveById(sideCarId string)
}

// cache is the Fanplane wrapper for go-control-plane cache. It uses Fanplane adapters to transform
// entities into to xDS Responses
type cache struct {
	envoyCache.SnapshotCache
	version atomic.Int32
}

func (c *cache) GetVersion() *atomic.Int32 {
	return &c.version
}

// NewCache returns a new valid Fanplane cache instance
func NewCache() *cache {
	result := cache{}
	result.SnapshotCache = envoyCache.NewSnapshotCache(true, NodeHash{}, log.StandardLogger())
	return &result
}

//Add object into Fanplane's cache if validated
func (c *cache) Add(obj fanplane.Object) (err error) {
	// Should recover from any adaptation error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf(
				"couldn't adapt %s for %s. %s",
				obj.GetObjectMeta().Name,
				obj.GetObjectKind().GroupVersionKind(),
				r,
			)
			log.Error(err)
		}
	}()

	cacheAdapter := adapter.NewAdapter(obj)
	rds, err := cacheAdapter.BuildRoutes()
	if err != nil {
		log.WithError(err).Errorln("couldn't parse routes.")
		return err
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
		return err
	}

	log.WithField("SicarId", obj.GetObjectMeta().Name).
		WithField("Type", obj.GetObjectKind().GroupVersionKind()).
		WithField("Version", fmt.Sprint(c.version.Load())).
		Infof("new version added in snapshot successfully.")

	return nil
}

// Remove forces invalidation of any snapshot that matches object metadata key name
func (c *cache) Remove(obj fanplane.Object) {
	c.SnapshotCache.ClearSnapshot(obj.GetObjectMeta().Name)
}

// Remove forces invalidation of any snapshot that matches sideCarId argument
func (c *cache) RemoveById(sideCarId string) {
	c.SnapshotCache.ClearSnapshot(sideCarId)
}
