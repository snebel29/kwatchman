package handler

import (
	"fmt"
	"testing"
)

func TestDiffHandlerUseSingletonStorage(t *testing.T) {
	s1 := newStorage()
	addr1 := fmt.Sprintf("%p", s1)
	s2 := newStorage()
	addr2 := fmt.Sprintf("%p", s2)
	if addr1 != addr2 {
		t.Errorf("addr should be the same however %s != %s", addr1, addr2)
	}
}
