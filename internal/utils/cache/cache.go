package cache

import (
	"errors"
)

var ErrCacheCapacity = errors.New("cache can't have 0 capacity")

type LRUCache struct {
	capacity int
	keys     []string
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

func (c *LRUCache) Get(key string) (string, bool) {
	for i, cKey := range c.keys {
		if key == cKey {
			c.keys = append(c.keys[:i], c.keys[i+1:]...)
			c.keys = append(c.keys, cKey)
			return c.itemsMap[key], true
		}
	}

	return "", false
}

func (c *LRUCache) Add(key, value string) bool {
	if _, exists := c.itemsMap[key]; exists {
		return false
	}

	if len(c.keys) == c.capacity {
		c.keys = c.keys[1:]
	}

	c.keys = append(c.keys, key)
	c.itemsMap[key] = value

	return true
}
