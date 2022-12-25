package hw04lrucache

import "sync"

type Cache[K comparable, V any] interface {
	Set(key K, value V) bool
	Get(key K) (V, bool)
	Clear()
}

type lruCache[K comparable, V any] struct {
	mu       sync.RWMutex
	capacity int
	queue    List[cacheItem[K, V]]
	items    map[K]*ListItem[cacheItem[K, V]]
}

type cacheItem[K comparable, V any] struct {
	key   K
	value V
}

func NewCache[K comparable, V any](capacity int) Cache[K, V] {
	return &lruCache[K, V]{
		capacity: capacity,
		queue:    NewList[cacheItem[K, V]](),
		items:    make(map[K]*ListItem[cacheItem[K, V]], capacity),
	}
}

func (l *lruCache[K, V]) Set(key K, value V) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	item, ok := l.items[key]
	if ok {
		item.Value.value = value
		l.queue.MoveToFront(item)
	} else {
		l.items[key] = l.queue.PushFront(cacheItem[K, V]{key: key, value: value})
	}

	if l.queue.Len() > l.capacity {
		backItem := l.queue.Back()
		l.queue.Remove(backItem)
		delete(l.items, backItem.Value.key)
	}

	return ok
}

func (l *lruCache[K, V]) Get(key K) (V, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var res V

	item, ok := l.items[key]
	if !ok {
		return res, false
	}
	l.queue.MoveToFront(item)
	return item.Value.value, true
}

func (l *lruCache[K, V]) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.queue = NewList[cacheItem[K, V]]()
	l.items = make(map[K]*ListItem[cacheItem[K, V]], l.capacity)
}

var _ Cache[string, int] = (*lruCache[string, int])(nil)
