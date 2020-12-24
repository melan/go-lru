package lru

import (
	"testing"
)

func TestListLRUCache(t *testing.T) {
	testLRUCache(t, NewListLRU)
}
