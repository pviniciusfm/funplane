package registry

import (
	"fmt"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	log "github.com/sirupsen/logrus"
	"github.frg.tech/cloud/fanplane/pkg/server"
	"go.uber.org/atomic"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/runtime"
	"github.com/ghodss/yaml"
)

type FileStore struct {
	Directory string
	version   atomic.Int32
	cache     cache.SnapshotCache
}


func NewFileRegistry(directoryPath string) (fileStore *FileStore, err error) {
	fileStore = &FileStore{}
	fileStore.cache = cache.NewSnapshotCache(true, server.Hasher{}, log.StandardLogger())
	err = fileStore.buildCache()
	return
}

func (store *FileStore) buildCache() (err error) {
	files, err := ioutil.ReadDir(store.Directory)
	if err != nil {
		return
	}

	for _, yamlFile := range files {
		snap := cache.NewSnapshot(fmt.Sprint(store.version.Inc()), endpoints, clusters, routes, listeners)
		rawYaml, err := ioutil.ReadFile(yamlFile.Name())
		header := runtime.TypeMeta{}

		err = yaml.Unmarshal(rawYaml, header)
		if err != nil {
			log.WithError(err).Errorf("could not parse %s.")
		}



		snapshotCache.SetSnapshot(node, snap)
	}

}
