package lru

import (
	"testing"
)

func TestMapLRUCache(t *testing.T) {
	testLRUCache(t, NewMapLRU)
}
