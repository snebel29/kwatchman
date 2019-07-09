package registry

import (
	"sync"
	log "github.com/sirupsen/logrus"
)

// Predefined registry identifiers
const (
	HANDLER   = "handler"
	RESOURCES = "resources" 
)

// ItemsRegistry holds a map of objects
type ItemsRegistry map[string]interface{}
var (
	lock = &sync.Mutex{}
	globalRegistry = map[string]ItemsRegistry{}
)

// Register items to the registry
func Register(registryName, itemName string, item interface{}) {
	log.Debugf("Registering item %#v with name %s into registry %s", item, itemName, registryName)
	lock.Lock()
	defer lock.Unlock()
	if globalRegistry[registryName] == nil {
		globalRegistry[registryName] = map[string]interface{}{}
	}
	globalRegistry[registryName][itemName] = item
}

// GetRegistry returns the global item registry
func GetRegistry(registryName string) (ItemsRegistry, bool) {
	log.Debugf("globalRegistry: %#v", globalRegistry)
	if r, exists := globalRegistry[registryName]; exists {
		return r, exists
	}
	return nil, false
}
