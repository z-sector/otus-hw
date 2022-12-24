package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu       sync.RWMutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	item, ok := l.items[key]
	if ok {
		l.queue.Remove(item)
	}

	if l.queue.Len() == l.capacity {
		backItem := l.queue.Back()
		l.queue.Remove(backItem)
		delete(l.items, backItem.Value.(cacheItem).key)
	}

	l.items[key] = l.queue.PushFront(cacheItem{key: key, value: value})

	return ok
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	item, ok := l.items[key]
	if !ok {
		return nil, false
	}
	l.queue.MoveToFront(item)
	return item.Value.(cacheItem).value, true
}

func (l *lruCache) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

var _ Cache = (*lruCache)(nil)
