package registry

import (
	"fmt"
	"github.frg.tech/cloud/fanplane/pkg/apis/fanplane"
	"github.frg.tech/cloud/fanplane/pkg/cache"
	"strings"
)

//Registry is a representation of a configuration registry where objects will be parsed and added in cache
type Registry interface {
	GetCache() cache.Cache
	GetConfig() *fanplane.Config
	GetRegistryType() Type
}

//Type is the current supported registry types
type Type string
const (
	KubeRegistry Type = "kube"
	FileRegistry      = "file"
)

// Convert the Level to a string. E.g. PanicLevel becomes "panic".
func (regType Type) String() string {
	switch regType {
	case KubeRegistry:
		return "kube"
	case FileRegistry:
		return "file"
	}
	return "unknown"
}

// ParseLevel takes a string level and returns the Logrus log level constant.
func ParseRegistryType(regType string) (Type, error) {
	switch strings.ToLower(regType) {
	case "kube":
		return KubeRegistry, nil
	case "file":
		return FileRegistry, nil
	}

	var l Type
	return l, fmt.Errorf("not a valid registry type: %q", regType)
}

// A constant exposing all registry types
var AllRegistryTypes = []Type{
	KubeRegistry,
	FileRegistry,
}
