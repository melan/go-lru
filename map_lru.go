package lru

type mapLRUItem struct {
	key             string
	value           interface{}
	hits            int
	morePopularNode *mapLRUItem
	lessPopularNode *mapLRUItem
}

type mapLRU struct {
	capacity       int
	cache          map[string]*mapLRUItem
	popularityTail *mapLRUItem
}

// NewMapLRU creates an instance of the LRU cache with a map as a backend
func NewMapLRU(capacity int) LRU {
	if capacity <= 0 {
		capacity = 1
	}

	return &mapLRU{
		capacity: capacity,
		cache:    make(map[string]*mapLRUItem, capacity),
	}
}

func (m *mapLRU) Get(key string) (found bool, value interface{}) {
	if item, ok := m.cache[key]; ok {
		item.hits++
		m.swap(item)
		return ok, item.value
	}

	return false, nil
}

func (m *mapLRU) Set(key string, value interface{}) {
	if item, ok := m.cache[key]; ok {
		item.hits++
		m.swap(item)
		return
	}

	if len(m.cache) == m.capacity {
		m.evict()
	}

	newItem := &mapLRUItem{
		key:             key,
		value:           value,
		hits:            1,
		morePopularNode: m.popularityTail,
	}
	m.cache[key] = newItem
	if m.popularityTail != nil {
		m.popularityTail.lessPopularNode = newItem
	}
	m.popularityTail = newItem
}

func (m *mapLRU) Size() int {
	return len(m.cache)
}

func (m *mapLRU) extractPopularityKeys() []string {
	keys := make([]string, 0, len(m.cache))

	for item := m.popularityTail; item != nil; item = item.morePopularNode {
		keys = append([]string{item.key}, keys...)
	}

	return keys
}

func (m *mapLRU) swap(item *mapLRUItem) {
	for {
		if item == nil || item.morePopularNode == nil || item.hits <= item.morePopularNode.hits {
			return
		}

		nextNode := item.morePopularNode
		nextNextNode := nextNode.morePopularNode

		nextNode.morePopularNode = item

		item.morePopularNode = nextNextNode
		if item.morePopularNode != nil {
			item.morePopularNode.lessPopularNode = item
		}

		nextNode.lessPopularNode = item.lessPopularNode
		if nextNode.lessPopularNode != nil {
			nextNode.lessPopularNode.morePopularNode = nextNode
		}

		item.lessPopularNode = nextNode

		if m.popularityTail == item {
			m.popularityTail = nextNode
		}
	}
}

func (m *mapLRU) evict() {
	if len(m.cache) < m.capacity || m.popularityTail == nil {
		return
	}

	item := m.popularityTail
	m.popularityTail = item.morePopularNode

	delete(m.cache, item.key)
}
