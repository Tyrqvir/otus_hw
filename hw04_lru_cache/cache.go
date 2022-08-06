package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value any) bool
	Get(key Key) (any, bool)
	Clear()
}

func (lc *lruCache) removeLastItem() {
	delete(lc.items, lc.queue.Back().Value.(cacheItem).key)
	lc.queue.Remove(lc.queue.Back())
}

func (lc *lruCache) Set(key Key, value any) bool {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	cached, isExist := lc.items[key]
	item := cacheItem{key: key, value: value}

	if isExist {
		cached.Value = item
		lc.queue.MoveToFront(cached)
	} else {
		lc.queue.PushFront(item)
		lc.items[key] = lc.queue.Front()
		if lc.isQueueMoreCapacity() {
			lc.removeLastItem()
		}
	}
	return isExist
}

func (lc *lruCache) Get(key Key) (any, bool) {
	lc.mu.RLock()
	defer lc.mu.RUnlock()

	cached, isExist := lc.items[key]
	if isExist {
		lc.queue.MoveToFront(cached)
		return lc.queue.Front().Value.(cacheItem).value, isExist
	}
	return nil, isExist
}

func (lc *lruCache) isQueueMoreCapacity() bool {
	return lc.queue.Len() > lc.capacity
}

func (lc *lruCache) Clear() {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	capacity := len(lc.items)
	lc.items = make(map[Key]*ListItem, capacity)
	lc.queue = NewList()
}

type lruCache struct {
	mu       sync.RWMutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value any
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
