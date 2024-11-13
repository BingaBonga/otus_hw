package hw04lrucache

import "sync"

type Key string

type Value struct {
	K Key
	V interface{}
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mutex    sync.RWMutex
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)

		itemValue, _ := item.Value.(*Value)
		itemValue.V = value
		return true
	}

	if len(c.items) >= c.capacity {
		itemValue, _ := c.queue.Back().Value.(*Value)

		delete(c.items, itemValue.K)
		c.queue.Remove(c.queue.Back())
	}

	c.items[key] = c.queue.PushFront(&Value{key, value})
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)

		itemValue, _ := item.Value.(*Value)
		return itemValue.V, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	defer c.mutex.Unlock()
	c.mutex.Lock()

	c.items = make(map[Key]*ListItem, c.capacity)
	c.queue = NewList()
}
