package lru

/*
	List-based implementation of the LRU cache. It increment hits for gets and for sets, meaning the most popular items will be on top of the list.

	This implementation isn't safe when accessed concurrently
*/

type listLRUItem struct {
	key   string
	value interface{}
	hits  int
}

type listLRU struct {
	cache []listLRUItem
}

// NewListLRU creates a new instance of the LRU cache
func NewListLRU(capacity int) LRU {
	if capacity <= 0 {
		capacity = 1
	}

	return &listLRU{
		cache: make([]listLRUItem, 0, capacity),
	}
}

func (l *listLRU) Get(key string) (found bool, value interface{}) {
	for i, item := range l.cache {
		if item.key == key {
			l.cache[i].hits++
			l.trySwap(i)

			return true, item.value
		}
	}
	return false, nil
}

func (l *listLRU) Set(key string, value interface{}) {
	for i, item := range l.cache {
		if item.key == key {
			l.cache[i].hits++
			l.trySwap(i)

			return
		}
	}

	if len(l.cache) == cap(l.cache) {
		l.cache[len(l.cache)-1].hits = 1
		l.cache[len(l.cache)-1].key = key
		l.cache[len(l.cache)-1].value = value
	} else {
		l.cache = append(l.cache, listLRUItem{hits: 1, key: key, value: value})
	}
}

func (l *listLRU) Size() int {
	return len(l.cache)
}

func (l *listLRU) trySwap(i int) {
	if i > 0 && l.cache[i].hits > l.cache[i-1].hits {
		prevItem := l.cache[i-1]
		l.cache[i-1].hits = l.cache[i].hits + 1
		l.cache[i-1].key = l.cache[i].key
		l.cache[i-1].value = l.cache[i].value

		l.cache[i].hits = prevItem.hits
		l.cache[i].key = prevItem.key
		l.cache[i].value = prevItem.value
	}
}
