package watcher

// Watcher interface
type Watcher interface {
	Run() error
	Shutdown()
}

// ResourceWatcher interface
type ResourceWatcher interface {
	Run() error
	Shutdown()
}
