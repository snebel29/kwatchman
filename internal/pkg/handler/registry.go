package handler

import (
	log "github.com/sirupsen/logrus"
)

type Registry map[string]Handler

var registry = Registry{}

// Register should be called only from init() on a handler package
func Register(name string, handler Handler) {
	log.Debugf("Registering handler %s", name)
	registry[name] = handler
}

func GetRegistry() Registry {
	return registry
}
