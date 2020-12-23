package lru

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListLRUCache(t *testing.T) {
	type cacheItem struct {
		key      string
		value    interface{}
		expected bool
	}

	tests := []struct {
		name     string
		capacity int
		items    []cacheItem
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
		},
		{
			name:     "Swap and Evictions",
			capacity: 2,
			items: []cacheItem{
				{
					key:      "a",
					value:    "b",
					expected: false,
				},
				{
					key:      "b",
					value:    "c",
					expected: true,
				},
				{
					key:      "b",
					value:    "c",
					expected: true,
				},
				{
					key:      "c",
					value:    "d",
					expected: true,
				},
			},
		},
	}

	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			cache := NewListLRU(tt.capacity)
			for _, item := range tt.items {
				cache.Set(item.key, item.value)
			}

			for _, item := range tt.items {
				gotFound, gotValue := cache.Get(item.key)

				if item.expected {
					assert.True(t, gotFound, "Item should be found")
					assert.Equal(t, item.value, gotValue)
				} else {
					assert.False(t, gotFound, "Item shouldn't be found")
					assert.Nil(t, gotValue)
				}
			}
		})
	}
}
