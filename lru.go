package lru

// LRU is an interface for different implementations of the LRU cache
type LRU interface {
	lruPopularityExtractor
	Get(key string) (bool, interface{})
	Set(key string, value interface{})
	Size() int
}

type lruPopularityExtractor interface {
	extractPopularityKeys() []string
}
