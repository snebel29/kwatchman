package watcher

type Watcher interface {
	Run() error
	Shutdown()
}

type ResourceWatcher interface {
	Run() error
	Shutdown()
}
