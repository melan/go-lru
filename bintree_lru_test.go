package lru

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBintreeLRUCache(t *testing.T) {
	testLRUCache(t, NewBintreeLRU)
}

type simplifiedBinTreeItem struct {
	left  string
	right string
}

func extractBinTreeItems(tip *bintreeLRUItem) map[string]simplifiedBinTreeItem {
	m := make(map[string]simplifiedBinTreeItem)
	if tip == nil {
		return m
	}

	var left string
	var right string

	if tip.left != nil {
		left = tip.left.key
	}

	if tip.right != nil {
		right = tip.right.key
	}

	m[tip.key] = simplifiedBinTreeItem{left: left, right: right}

	leftMap := extractBinTreeItems(tip.left)
	for k, v := range leftMap {
		m[k] = v
	}

	rightMap := extractBinTreeItems(tip.right)
	for k, v := range rightMap {
		m[k] = v
	}

	return m
}

func TestBintreeLRUCache_Set(t *testing.T) {
	tests := []struct {
		name              string
		items             []string
		capacity          int
		wantSize          int
		wantItemsPriority []string
		wantItemsTree     map[string]simplifiedBinTreeItem
	}{
		{
			name:     "not at capacity",
			items:    []string{"a", "b", "c", "d", "e"},
			capacity: 10,
			wantSize: 5,
			wantItemsPriority: []string{
				"a", "b", "c", "d", "e"},
			wantItemsTree: map[string]simplifiedBinTreeItem{
				"a": {left: "", right: ""},
				"b": {left: "a", right: "d"},
				"c": {left: "", right: ""},
				"d": {left: "c", right: "e"},
				"e": {left: "", right: ""},
			},
		},
		{
			name:     "not at capacity, reversed",
			items:    []string{"e", "d", "c", "b", "a"},
			capacity: 10,
			wantSize: 5,
			wantItemsPriority: []string{
				"e", "d", "c", "b", "a"},
			wantItemsTree: map[string]simplifiedBinTreeItem{
				"a": {left: "", right: ""},
				"b": {left: "a", right: "c"},
				"c": {left: "", right: ""},
				"d": {left: "b", right: "e"},
				"e": {left: "", right: ""},
			},
		},
		{
			name:     "not at capacity, duplicates",
			items:    []string{"e", "e", "e", "e", "e"},
			capacity: 10,
			wantSize: 1,
			wantItemsPriority: []string{
				"e"},
			wantItemsTree: map[string]simplifiedBinTreeItem{
				"e": {left: "", right: ""},
			},
		},
		{
			name:     "not at capacity, longer",
			items:    []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o"},
			capacity: 20,
			wantSize: 15,
			wantItemsPriority: []string{
				"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o"},
			wantItemsTree: map[string]simplifiedBinTreeItem{
				"a": {left: "", right: ""},
				"b": {left: "a", right: "c"},
				"c": {left: "", right: ""},
				"d": {left: "b", right: "f"},
				"e": {left: "", right: ""},
				"f": {left: "e", right: "g"},
				"g": {left: "", right: ""},
				"h": {left: "d", right: "l"},
				"i": {left: "", right: ""},
				"j": {left: "i", right: "k"},
				"k": {left: "", right: ""},
				"l": {left: "j", right: "n"},
				"m": {left: "", right: ""},
				"n": {left: "m", right: "o"},
				"o": {left: "", right: ""},
			},
		},
		{
			name:     "evictoins",
			items:    []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o"},
			capacity: 10,
			wantSize: 10,
			wantItemsPriority: []string{
				"a", "b", "c", "d", "e", "f", "g", "h", "i", "o"},
			wantItemsTree: map[string]simplifiedBinTreeItem{
				"a": {left: "", right: ""},
				"b": {left: "a", right: "c"},
				"c": {left: "", right: ""},
				"d": {left: "b", right: "h"},
				"e": {left: "", right: ""},
				"f": {left: "e", right: "g"},
				"g": {left: "", right: ""},
				"h": {left: "f", right: "i"},
				"i": {left: "", right: "o"},
				"o": {left: "", right: ""},
			},
		},
		{
			name:     "evictoins with duplicates",
			items:    strings.Split("jaskldfhcweoichpqwoiehcmkamshjcfnioqhwecfionhqpwiehfluvnhwrbiuvhbnsihdfbviavwheoifanwioefhcqhuierhvboaiuwehcnofiquwhefoihgbahvoimacmjfoniahjwoeihvblaushdlfajkvshldlvjnkshcmiuehbghvlaksndfmvzxmnhfvuiahberoigupvhqwbeghpbqvuweyrpinvqwkjsbdvjkdzhalnviuwevoybuiwehcfmoiquwhenivubqhwpeiufcmhlskduhlfaishdlfabjkhsdflivhaslkjdfhlbaiushvlnviufhalwuiehfkjshdflbiuvahwleuifhiwubhvajkshdf", ""),
			capacity: 10,
			wantSize: 10,
			wantItemsPriority: []string{
				"h", "f", "a", "l", "s", "j", "k", "c", "m", "d"},
			wantItemsTree: map[string]simplifiedBinTreeItem{
				"a": {left: "", right: ""},
				"c": {left: "a", right: "d"},
				"d": {left: "", right: ""},
				"f": {left: "c", right: "j"},
				"h": {left: "", right: ""},
				"j": {left: "h", right: ""},
				"k": {left: "f", right: "m"},
				"l": {left: "", right: ""},
				"m": {left: "l", right: "s"},
				"s": {left: "", right: ""},
			},
		},
	}

	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			cache := NewBintreeLRU(tt.capacity).(*bintreeLRU)
			for _, key := range tt.items {
				cache.Set(key, "")
			}

			assert.Equal(t, tt.wantSize, cache.Size())
			gotPopularityList := cache.extractPopularityKeys()

			assert.Equal(t, tt.wantItemsPriority, gotPopularityList)

			gotTreeStruct := extractBinTreeItems(cache.tip)
			assert.Equal(t, tt.wantItemsTree, gotTreeStruct)
		})
	}
}

func TestBintreeLRUCache_findBiggestInSubTree(t *testing.T) {

	tests := []struct {
		name        string
		items       []string
		capacity    int
		wantElement string
	}{
		{
			name:        "direct left",
			items:       strings.Split("abcde", ""),
			capacity:    10,
			wantElement: "a",
		},
		{
			name:        "deep right",
			items:       strings.Split("edcba", ""),
			capacity:    10,
			wantElement: "c",
		},
	}

	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			cache := NewBintreeLRU(tt.capacity).(*bintreeLRU)
			for _, key := range tt.items {
				cache.Set(key, "")
			}

			gotBiggest := cache.findBiggestInSubtree(cache.tip.left)
			assert.Equal(t, tt.wantElement, gotBiggest.key)
		})
	}
}

func TestBintreeLRUCache_evict(t *testing.T) {
	tests := []struct {
		name              string
		items             []string
		capacity          int
		wantSize          int
		wantItemsPriority []string
		wantItemsTree     map[string]simplifiedBinTreeItem
	}{
		{
			name:              "evict tip 1",
			items:             strings.Split("abcdeacde", ""),
			capacity:          5,
			wantSize:          4,
			wantItemsPriority: strings.Split("acde", ""),
			wantItemsTree: map[string]simplifiedBinTreeItem{
				"a": {left: "", right: "c"},
				"c": {left: "", right: ""},
				"d": {left: "a", right: "e"},
				"e": {left: "", right: ""},
			},
		},
		{
			name:              "evict tip 2",
			items:             strings.Split("edcbaabce", ""),
			capacity:          5,
			wantSize:          4,
			wantItemsPriority: strings.Split("abce", ""),
			wantItemsTree: map[string]simplifiedBinTreeItem{
				"a": {left: "", right: ""},
				"b": {left: "a", right: ""},
				"c": {left: "b", right: "e"},
				"e": {left: "", right: ""},
			},
		},
		{
			name:              "evict node 1",
			items:             strings.Split("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwyz", ""),
			capacity:          26,
			wantSize:          25,
			wantItemsPriority: strings.Split("abcdefghijklmnopqrstuvwyz", ""),
			wantItemsTree: map[string]simplifiedBinTreeItem{
				"a": {left: "", right: ""},
				"b": {left: "a", right: "c"},
				"c": {left: "", right: ""},
				"d": {left: "b", right: "f"},
				"e": {left: "", right: ""},
				"f": {left: "e", right: "g"},
				"g": {left: "", right: ""},
				"h": {left: "d", right: "l"},
				"i": {left: "", right: ""},
				"j": {left: "i", right: "k"},
				"k": {left: "", right: ""},
				"l": {left: "j", right: "n"},
				"m": {left: "", right: ""},
				"n": {left: "m", right: "o"},
				"o": {left: "", right: ""},
				"p": {left: "h", right: "t"},
				"q": {left: "", right: ""},
				"r": {left: "q", right: "s"},
				"s": {left: "", right: ""},
				"t": {left: "r", right: "w"},
				"u": {left: "", right: ""},
				"v": {left: "u", right: ""},
				"w": {left: "v", right: "y"},
				"y": {left: "", right: "z"},
				"z": {left: "", right: ""},
			},
		},
		{
			name:              "evict node 2",
			items:             strings.Split("abcdefghijklmnopqrstuvwxyzabcdefgijklmnopqrstuvwxyz", ""),
			capacity:          26,
			wantSize:          25,
			wantItemsPriority: strings.Split("abcdefgijklmnopqrstuvwxyz", ""),
			wantItemsTree: map[string]simplifiedBinTreeItem{
				"a": {left: "", right: ""},
				"b": {left: "a", right: "c"},
				"c": {left: "", right: ""},
				"d": {left: "b", right: "f"},
				"e": {left: "", right: ""},
				"f": {left: "e", right: ""},
				"g": {left: "d", right: "l"},
				"i": {left: "", right: ""},
				"j": {left: "i", right: "k"},
				"k": {left: "", right: ""},
				"l": {left: "j", right: "n"},
				"m": {left: "", right: ""},
				"n": {left: "m", right: "o"},
				"o": {left: "", right: ""},
				"p": {left: "g", right: "t"},
				"q": {left: "", right: ""},
				"r": {left: "q", right: "s"},
				"s": {left: "", right: ""},
				"t": {left: "r", right: "x"},
				"u": {left: "", right: ""},
				"v": {left: "u", right: "w"},
				"w": {left: "", right: ""},
				"x": {left: "v", right: "y"},
				"y": {left: "", right: "z"},
				"z": {left: "", right: ""},
			},
		},
	}

	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			cache := NewBintreeLRU(tt.capacity).(*bintreeLRU)
			for _, key := range tt.items {
				cache.Set(key, "")
			}

			cache.evict()

			assert.Equal(t, tt.wantSize, cache.Size())
			gotPopularityList := cache.extractPopularityKeys()

			assert.Equal(t, tt.wantItemsPriority, gotPopularityList)

			gotTreeStruct := extractBinTreeItems(cache.tip)
			assert.Equal(t, tt.wantItemsTree, gotTreeStruct)
		})
	}
}
