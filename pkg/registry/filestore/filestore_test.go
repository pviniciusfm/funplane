package filestore

import (
	"github.com/stretchr/testify/assert"
	"github.frg.tech/cloud/fanplane/pkg/apis/fanplane"
	"github.frg.tech/cloud/fanplane/pkg/cache"
	"github.frg.tech/cloud/fanplane/pkg/registry"
	"go.uber.org/atomic"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

//mockCache is an mock implementation of server.ServerCache that helps match expected objects in the cacheStore
type mockCache struct {
	cache.Cache
	version      atomic.Int32
	expectations map[string]int32
	cacheStore   map[string]int32
}

func (mock *mockCache) GetVersion() *atomic.Int32 {
	return &mock.version
}

func (mock *mockCache) Add(obj fanplane.Object) (err error) {
	mock.cacheStore[obj.GetObjectMeta().Name] = mock.version.Inc()
	return nil
}

func (mock *mockCache) Remove(obj fanplane.Object) {
	delete(mock.cacheStore, obj.GetObjectMeta().Name)
}

func (mock *mockCache) RemoveById(sideCarId string) {
	delete(mock.cacheStore, sideCarId)
}

func NewMockCache() *mockCache {
	return &mockCache{cacheStore: map[string]int32{}}
}

func TestAddRegistry(t *testing.T) {
	cache := NewMockCache()
	fileRegistryPath, err := ioutil.TempDir("", "fanplane")
	if err != nil {
		t.Fail()
	}
	defer os.Remove(fileRegistryPath)

	config := &fanplane.Config{RegistryType: registry.FileRegistry, RegistryDirectory: fileRegistryPath}
	fileRegistry, err := NewFileRegistry(config, cache)
	if err != nil {
		t.FailNow()
	}

	err = fileRegistry.Add("../testdata/envoy.yaml")
	if err != nil {
		t.FailNow()
	}

	assert.Equal(t, int32(1), cache.version.Load())
	assert.EqualValues(t, 1, len(cache.cacheStore))
	assert.Contains(t, cache.cacheStore, "falcon")
}

func TestRemoveRegistry(t *testing.T) {
	cache := NewMockCache()
	cache.cacheStore["falcon"] = cache.GetVersion().Inc()

	fileRegistryPath, err := ioutil.TempDir("", "fanplane")
	if err != nil {
		t.Fail()
	}
	defer os.Remove(fileRegistryPath)

	config := &fanplane.Config{RegistryType: registry.FileRegistry, RegistryDirectory: fileRegistryPath}
	fileRegistry, err := NewFileRegistry(config, cache)
	fileRegistry.files["/mock/test/data"] = "falcon"
	checkErr(t, err)

	err = fileRegistry.Remove("/mock/test/data")
	checkErr(t, err)

	assert.Equal(t, 0, len(fileRegistry.files))
	//FUTURE: support fetch in Fanplane wrapped server cache interface
}

func TestBuildCache(t *testing.T) {
	cache := NewMockCache()
	fileRegistryPath, err := ioutil.TempDir("", "fanplane")
	if err != nil {
		t.Fail()
	}
	defer os.Remove(fileRegistryPath)
	copy(t, "../testdata/envoy.yaml", path.Join(fileRegistryPath, "envoy.yaml"))

	config := &fanplane.Config{RegistryType: registry.FileRegistry, RegistryDirectory: fileRegistryPath}
	fileRegistry, err := NewFileRegistry(config, cache)

	assert.Nil(t, err)
	assert.NotNil(t, fileRegistry)
	assert.Equal(t, int32(1), cache.version.Load())
}

func copy(t *testing.T, src string, dst string) {
	// Read all content of src to data
	data, err := ioutil.ReadFile(src)
	checkErr(t, err)

	// Write data to dst
	err = ioutil.WriteFile(dst, data, 0644)
	checkErr(t, err)
}

func checkErr(t *testing.T, e error) {
	if e != nil {
		t.FailNow()
	}
}
