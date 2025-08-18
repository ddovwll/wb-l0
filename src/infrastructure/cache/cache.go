package cache

import (
	"container/list"
	"errors"
	"sync"
)

type Item struct {
	Key   string
	Value any
}

type LRUCache struct {
	capacity int
	items    map[string]*list.Element
	list     *list.List
	mutex    sync.Mutex
}

func NewLRUCache(capacity int) (*LRUCache, error) {
	if capacity <= 0 {
		return nil, errors.New("capacity must be greater than zero")
	}

	return &LRUCache{
			capacity: capacity,
			items:    make(map[string]*list.Element),
			list:     list.New(),
		},
		nil
}

func (cache *LRUCache) Get(key string) any {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	item, ok := cache.items[key]
	if !ok {
		return nil
	}

	cache.list.MoveToFront(item)
	return item.Value.(*Item).Value
}

func (cache *LRUCache) Set(key string, value any) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	if element, exists := cache.items[key]; exists {
		cache.list.MoveToFront(element)
		element.Value.(*Item).Value = value
	}

	if cache.list.Len() == cache.capacity {
		cache.remove()
	}

	item := &Item{
		Key:   key,
		Value: value,
	}

	element := cache.list.PushFront(item)
	cache.items[key] = element
}

func (cache *LRUCache) remove() {
	if element := cache.list.Back(); element != nil {
		item := cache.list.Remove(element).(*Item)
		delete(cache.items, item.Key)
	}
}
