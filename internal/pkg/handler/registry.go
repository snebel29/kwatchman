package handler

import (
	log "github.com/sirupsen/logrus"
)

// This is the handler registry
var Registry = map[string]Handler{}

// Register should be called from init() on a handler package
func Register(name string, handler Handler) {
	log.Debugf("Registering handler %s", name)
	Registry[name] = handler
}
