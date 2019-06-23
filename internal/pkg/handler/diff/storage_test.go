package diff

import (
	"reflect"
	"testing"
)

func TestStorage_Add(t *testing.T) {
	s := newStorage()
	key := "a"
	value := []byte("a_value")

	s.Add(key, value)
	if !s.Has(key) {
		t.Errorf("Storage should have key %s", key)
	}

	v, ok := s.Get(key)
	if !reflect.DeepEqual(v, value) {
		t.Errorf("returned value should have beem %#v", value)
	}
	if ok != true {
		t.Error("ok should be true")
	}

	s.Delete(key)
	if s.Has(key) {
		t.Errorf("Storage should NOT have key %s", key)
	}
}
