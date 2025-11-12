package service

import (
	"container/list"
	"sync"

	"github.com/sonni-a/wb-service/internal/models"
)

type Cache interface {
	Get(key string) (*models.Order, bool)
	Set(key string, value *models.Order)
	Delete(key string)
	Clear()
}

type MemoryCache struct {
	mu      sync.RWMutex
	items   map[string]*list.Element
	order   *list.List
	maxSize int
}

type cacheEntry struct {
	key   string
	value *models.Order
}

func NewMemoryCache(maxSize int) Cache {
	return &MemoryCache{
		items:   make(map[string]*list.Element),
		order:   list.New(),
		maxSize: maxSize,
	}
}

func (c *MemoryCache) Get(key string) (*models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if el, ok := c.items[key]; ok {
		c.mu.RUnlock()
		c.mu.Lock()
		c.order.MoveToFront(el)
		c.mu.Unlock()
		c.mu.RLock()
		return el.Value.(*cacheEntry).value, true
	}
	return nil, false
}

func (c *MemoryCache) Set(key string, value *models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if el, ok := c.items[key]; ok {
		el.Value.(*cacheEntry).value = value
		c.order.MoveToFront(el)
		return
	}

	entry := &cacheEntry{key, value}
	el := c.order.PushFront(entry)
	c.items[key] = el

	if c.order.Len() > c.maxSize {
		oldest := c.order.Back()
		if oldest != nil {
			c.order.Remove(oldest)
			delete(c.items, oldest.Value.(*cacheEntry).key)
		}
	}
}

func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if el, ok := c.items[key]; ok {
		c.order.Remove(el)
		delete(c.items, key)
	}
}

func (c *MemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*list.Element)
	c.order.Init()
}
