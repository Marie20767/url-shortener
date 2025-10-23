package cache

import (
	"errors"
)

var (
	ErrCacheMiss               = errors.New("cache miss")
	ErrCacheEntryAlreadyExists = errors.New("entry with this key already exists in cache")
	ErrCacheCapacity           = errors.New("cache can't have 0 capacity")
)

type LRUCache struct {
	capacity int
	keys []string
	itemsMap map[string]string
}

func New(capacity int) (*LRUCache, error) {
	if capacity == 0 {
		return nil, ErrCacheCapacity
	}

	return &LRUCache{
		capacity: capacity,
		keys:     []string{},
		itemsMap: map[string]string{},
	}, nil
}

func (c *LRUCache) Get(key string) (string, error) {
	for i, cKey := range c.keys {
		if key == cKey {
			c.keys = append(c.keys[:i], c.keys[i+1:]...)
			c.keys = append(c.keys, cKey)
			return c.itemsMap[key], nil
		}
	}

	return "", ErrCacheMiss
}

func (c *LRUCache) Add(key string, value string) error {
	if _, exists := c.itemsMap[key]; exists {
		return ErrCacheEntryAlreadyExists
	}

	if len(c.keys) == c.capacity {
		c.keys = c.keys[1:]
	}

	c.keys = append(c.keys, key)
	c.itemsMap[key] = value

	return nil
}
