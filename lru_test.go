package lru

import (
	"fmt"
	"strings"
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
		name              string
		capacity          int
		items             []string
		wantSize          int
		wantItemsPriority []string
		wantItems         []string
	}{
		{
			name:              "Single item",
			capacity:          2,
			items:             []string{"a"},
			wantItemsPriority: []string{"a"},
			wantItems:         []string{"a"},
			wantSize:          1,
		},
		{
			name:              "Zero capacity",
			capacity:          0,
			items:             []string{"a", "c"},
			wantItemsPriority: []string{"c"},
			wantItems:         []string{"c"},
			wantSize:          1,
		},
		{
			name:              "Exact cache capacity",
			capacity:          2,
			items:             []string{"a", "b"},
			wantItemsPriority: []string{"a", "b"},
			wantItems:         []string{"a", "b"},
			wantSize:          2,
		},
		{
			name:              "Duplicates",
			capacity:          2,
			items:             strings.Split("abaa", ""),
			wantItemsPriority: []string{"a", "b"},
			wantItems:         []string{"a", "b"},
			wantSize:          2,
		},
		{
			name:              "Evictions",
			capacity:          2,
			items:             strings.Split("abca", ""),
			wantItemsPriority: strings.Split("ac", ""),
			wantItems:         strings.Split("ac", ""),
			wantSize:          2,
		},
		{
			name:              "Duplicates, Swaps and Evictions",
			capacity:          4,
			items:             strings.Split("abbbcazccczzbddzzzcddddcdcbbeeeeeeeeedccc", ""),
			wantItemsPriority: strings.Split("cedz", ""),
			wantItems:         strings.Split("cedz", ""),
			wantSize:          4,
		},
		{
			name:              "evictoins with duplicates",
			items:             strings.Split("jaskldfhcweoichpqwoiehcmkamshjcfnioqhwecfionhqpwiehfluvnhwrbiuvhbnsihdfbviavwheoifanwioefhcqhuierhvboaiuwehcnofiquwhefoihgbahvoimacmjfoniahjwoeihvblaushdlfajkvshldlvjnkshcmiuehbghvlaksndfmvzxmnhfvuiahberoigupvhqwbeghpbqvuweyrpinvqwkjsbdvjkdzhalnviuwevoybuiwehcfmoiquwhenivubqhwpeiufcmhlskduhlfaishdlfabjkhsdflivhaslkjdfhlbaiushvlnviufhalwuiehfkjshdflbiuvahwleuifhiwubhvajkshdf", ""),
			capacity:          10,
			wantSize:          10,
			wantItemsPriority: strings.Split("hfalsjkcmd", ""),
		},
	}

	testValue := "some value"
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			cache := newLRU(tt.capacity)
			for _, key := range tt.items {
				cache.Set(key, testValue)
			}

			assert.Equal(t, tt.wantSize, cache.Size())

			gotPopularityList := cache.extractPopularityKeys()
			assert.Equal(t, tt.wantItemsPriority, gotPopularityList, "Popularity list doesn't match")

			for _, key := range tt.wantItems {
				gotFound, gotValue := cache.Get(key)

				assert.True(t, gotFound, fmt.Sprintf("Item %q should be found", key))
				assert.Equal(t, testValue, gotValue, fmt.Sprintf("Value of %q mismatches", key))
			}
		})
	}
}
