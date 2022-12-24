package hw04lrucache

import (
	"sync"
)

type Key string

type Cache[T any] interface {
	Set(key Key, value T) bool
	Get(key Key) (T, bool)
	Clear()
}

type lruCache[T any] struct {
	mu       sync.RWMutex
	capacity int
	queue    List[cacheItem[T]]
	items    map[Key]*ListItem[cacheItem[T]]
}

type cacheItem[T any] struct {
	key   Key
	value T
}

func NewCache[T any](capacity int) Cache[T] {
	return &lruCache[T]{
		capacity: capacity,
		queue:    NewList[cacheItem[T]](),
		items:    make(map[Key]*ListItem[cacheItem[T]], capacity),
	}
}

func (l *lruCache[T]) Set(key Key, value T) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	item, ok := l.items[key]
	if ok {
		l.queue.Remove(item)
	}

	if l.queue.Len() == l.capacity {
		backItem := l.queue.Back()
		l.queue.Remove(backItem)
		delete(l.items, backItem.Value.key)
	}

	l.items[key] = l.queue.PushFront(cacheItem[T]{key: key, value: value})

	return ok
}

func (l *lruCache[T]) Get(key Key) (T, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var res T

	item, ok := l.items[key]
	if !ok {
		return res, false
	}
	l.queue.MoveToFront(item)
	return item.Value.value, true
}

func (l *lruCache[T]) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.queue = NewList[cacheItem[T]]()
	l.items = make(map[Key]*ListItem[cacheItem[T]], l.capacity)
}

var _ Cache[int] = (*lruCache[int])(nil)
