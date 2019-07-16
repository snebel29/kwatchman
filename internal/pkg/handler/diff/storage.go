package diff

import "sync"

type storage struct {
	sync.RWMutex
	repository map[string][]byte
}

// Return a new storage object, every diff handler should have its own
// so singleton pattern was finally not necessary
func newStorage() *storage {
	r := make(map[string][]byte)
	return &storage{
		repository: r,
	}
}

func (s *storage) Add(key string, value []byte) {
	s.Lock()
	defer s.Unlock()
	s.repository[key] = value
}

func (s *storage) Delete(key string) {
	s.Lock()
	defer s.Unlock()
	delete(s.repository, key)
}

func (s *storage) Get(key string) ([]byte, bool) {
	s.RLock()
	defer s.RUnlock()
	v, ok := s.repository[key]
	return v, ok
}

func (s *storage) Has(key string) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.repository[key]
	return ok
}
