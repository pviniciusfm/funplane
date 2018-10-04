package filestore

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/hashicorp/go-multierror"
	log "github.com/sirupsen/logrus"
	"github.frg.tech/cloud/fanplane/pkg/apis/fanplane"
	"github.frg.tech/cloud/fanplane/pkg/apis/fanplane/v1alpha1"
	"github.frg.tech/cloud/fanplane/pkg/cache"
	"github.frg.tech/cloud/fanplane/pkg/registry"
	"os"
	"path/filepath"
	"sync"
)

//fileRegistry is a representation of a config registry in filesystem
type fileRegistry struct {
	cache  cache.Cache
	config *fanplane.Config
	files  map[string]string
	mu     sync.Mutex
}

func (store *fileRegistry) GetCache() cache.Cache {
	return store.cache
}

func (store *fileRegistry) GetConfig() *fanplane.Config {
	return store.config
}

func (store *fileRegistry) GetRegistryType() registry.Type {
	return registry.FileRegistry
}

//NewFileRegistry return new instance of FileRegistry
func NewFileRegistry(config *fanplane.Config, cache cache.Cache) (newRegistry *fileRegistry, err error) {
	if config.RegistryDirectory == "" {
		return nil, fmt.Errorf("a directory path is required")
	}

	if _, err = os.Stat(config.RegistryDirectory); os.IsNotExist(err) {
		return nil, fmt.Errorf("a directory path %q does not exist", config.RegistryDirectory)
	}

	newRegistry = &fileRegistry{config: config, cache: cache, files: map[string]string{}}
	err = newRegistry.buildCache()
	return
}

// buildCache walks through registry folder adding objects into server cacheStore
func (store *fileRegistry) buildCache() (error) {
	var errors error
	var files []string

	err := filepath.Walk(store.config.RegistryDirectory, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return multierror.Append(err, fmt.Errorf("couldn't search in directory. %s", err))
	}

	for _, yamlFile := range files {
		log.Debugf("yaml file detected %s ", yamlFile)

		if yamlFile != "" {
			err = store.Add(yamlFile)
			if err != nil {
				continue
			}
		}

		return errors
	}

	return nil
}

//Add method is used to add a new Fanplane object into registry and server cacheStore
func (store *fileRegistry) Add(filePath string) (err error) {
	kind, err := v1alpha1.LoadFanplaneKind(filePath)
	if err != nil || !v1alpha1.IsValidFanplaneObject(kind.Kind) {
		err = fmt.Errorf("couldn't add %s in the registry %s", filePath, err)
		log.Error(err)
		return
	}

	store.mu.Lock()
	store.files[filePath] = kind.Name
	store.mu.Unlock()

	store.cache.Add(kind)
	return
}

//Remove method is used to remove the registry files from map and cacheStore
func (store *fileRegistry) Remove(filePath string) (err error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	metadata := store.files[filePath]

	if metadata != "" {
		store.cache.RemoveById(store.files[filePath])
		delete(store.files, filePath)
		log.Infof("removed %s file with metadata %s ", filePath, metadata)
	}

	return
}

//StartWatch starts the monitor loop receivent events from file system in case of creating, updating and removing files
func (store *fileRegistry) StartWatch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				log.Debug("event:", event)
				if event.Op&fsnotify.Create == fsnotify.Create {
					store.Add(event.Name)
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					store.Add(event.Name)
				}

				if event.Op&fsnotify.Remove == fsnotify.Remove {
					store.Remove(event.Name)
				}
			case watchErr, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Error("error:", watchErr)
			}
		}
	}()

	log.Infof("fileRegistry start watching %s", store.config.RegistryDirectory)

	err = watcher.Add(store.config.RegistryDirectory)

	if err != nil {
		log.Fatal(err)
	}
	<-done
}
