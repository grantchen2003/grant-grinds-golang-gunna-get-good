package main

import "sync"

type ConcurrentSafeMap[K comparable, V any] struct {
	mutex sync.RWMutex
	data  map[K]V
}

func NewConcurrentSafeMap[K comparable, V any]() *ConcurrentSafeMap[K, V] {
	return &ConcurrentSafeMap[K, V]{
		mutex: sync.RWMutex{},
		data:  make(map[K]V),
	}
}

func (c *ConcurrentSafeMap[K, V]) Set(key K, value V) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = value
}

func (c *ConcurrentSafeMap[K, V]) Get(key K) (V, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	value, exists := c.data[key]
	return value, exists
}

func (c *ConcurrentSafeMap[K, V]) Delete(key K) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.data, key)
}

func (c *ConcurrentSafeMap[K, V]) Exists(key K) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	_, exists := c.data[key]
	return exists
}

func (c *ConcurrentSafeMap[K, V]) Len() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return len(c.data)
}

func (c *ConcurrentSafeMap[K, V]) Keys() []K {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	keys := make([]K, 0, len(c.data))
	for key := range c.data {
		keys = append(keys, key)
	}
	return keys
}
