package lru

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testLRUCache(t *testing.T, newLRU func(capacity int) LRU) {
	type cacheItem struct {
		key      string
		value    interface{}
		expected bool
	}

	tests := []struct {
		name     string
		capacity int
		items    []cacheItem
		wantSize int
	}{
		{
			name:     "Single item",
			capacity: 2,
			items: []cacheItem{
				{
					key:      "a",
					value:    "b",
					expected: true,
				},
			},
			wantSize: 1,
		},
		{
			name:     "Zero capacity",
			capacity: 0,
			items: []cacheItem{
				{
					key:      "a",
					value:    "b",
					expected: false,
				},
				{
					key:      "c",
					value:    "d",
					expected: true,
				},
			},
			wantSize: 1,
		},
		{
			name:     "Exact cache capacity",
			capacity: 2,
			items: []cacheItem{
				{
					key:      "a",
					value:    "b",
					expected: true,
				},
				{
					key:      "b",
					value:    "c",
					expected: true,
				},
			},
			wantSize: 2,
		},
		{
			name:     "Duplicates",
			capacity: 2,
			items: []cacheItem{
				{
					key:      "a",
					value:    "b",
					expected: true,
				},
				{
					key:      "b",
					value:    "c",
					expected: true,
				},
				{
					key:      "a",
					value:    "b",
					expected: true,
				},
				{
					key:      "a",
					value:    "b",
					expected: true,
				},
			},
			wantSize: 2,
		},
		{
			name:     "Evictions",
			capacity: 2,
			items: []cacheItem{
				{
					key:      "a",
					value:    "b",
					expected: true,
				},
				{
					key:      "b",
					value:    "c",
					expected: false,
				},
				{
					key:      "c",
					value:    "d",
					expected: true,
				},
				{
					key:      "a",
					value:    "b",
					expected: true,
				},
			},
			wantSize: 2,
		},
		{
			name:     "Duplicates, Swaps and Evictions",
			capacity: 4,
			items: []cacheItem{
				{key: "a", value: " ", expected: false},
				{key: "b", value: " ", expected: false},
				{key: "b", value: " ", expected: false},
				{key: "b", value: " ", expected: false},
				{key: "c", value: " ", expected: true},
				{key: "a", value: " ", expected: false},
				{key: "z", value: " ", expected: true},
				{key: "c", value: " ", expected: true},
				{key: "c", value: " ", expected: true},
				{key: "c", value: " ", expected: true},
				{key: "z", value: " ", expected: true},
				{key: "z", value: " ", expected: true},
				{key: "b", value: " ", expected: false},
				{key: "d", value: " ", expected: true},
				{key: "d", value: " ", expected: true},
				{key: "z", value: " ", expected: true},
				{key: "z", value: " ", expected: true},
				{key: "z", value: " ", expected: true},
				{key: "c", value: " ", expected: true},
				{key: "d", value: " ", expected: true},
				{key: "d", value: " ", expected: true},
				{key: "d", value: " ", expected: true},
				{key: "d", value: " ", expected: true},
				{key: "c", value: " ", expected: true},
				{key: "d", value: " ", expected: true},
				{key: "c", value: " ", expected: true},
				{key: "b", value: " ", expected: false},
				{key: "b", value: " ", expected: false},
				{key: "e", value: " ", expected: true},
				{key: "e", value: " ", expected: true},
				{key: "e", value: " ", expected: true},
				{key: "e", value: " ", expected: true},
				{key: "e", value: " ", expected: true},
				{key: "e", value: " ", expected: true},
				{key: "e", value: " ", expected: true},
				{key: "e", value: " ", expected: true},
				{key: "e", value: " ", expected: true},
				{key: "d", value: " ", expected: true},
				{key: "c", value: " ", expected: true},
				{key: "c", value: " ", expected: true},
				{key: "c", value: " ", expected: true},
			},
			wantSize: 4,
		},
	}

	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			cache := newLRU(tt.capacity)
			for _, item := range tt.items {
				cache.Set(item.key, item.value)
			}

			assert.Equal(t, tt.wantSize, cache.Size())

			for _, item := range tt.items {
				gotFound, gotValue := cache.Get(item.key)

				if item.expected {
					assert.True(t, gotFound, fmt.Sprintf("Item %q should be found", item.key))
					assert.Equal(t, item.value, gotValue, fmt.Sprintf("Value of %q mismatches", item.key))
				} else {
					assert.False(t, gotFound, fmt.Sprintf("Item %q shouldn't be found", item.key))
					assert.Nil(t, gotValue, fmt.Sprintf("Value of %q isn't nil", item.key))
				}
			}
		})
	}
}
