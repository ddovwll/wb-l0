package cache

import (
	"container/list"
	"errors"
	"sync"
)

type Item[Key comparable, Value any] struct {
	key   Key
	value Value
}

type LRUCache[Key comparable, Value any] struct {
	capacity int
	items    map[Key]*list.Element
	list     *list.List
	mutex    sync.Mutex
}

func NewLRUCache[Key comparable, Value any](capacity int) (*LRUCache[Key, Value], error) {
	if capacity <= 0 {
		return nil, errors.New("capacity must be greater than zero")
	}

	return &LRUCache[Key, Value]{
			capacity: capacity,
			items:    make(map[Key]*list.Element),
			list:     list.New(),
		},
		nil
}

func (cache *LRUCache[Key, Value]) Get(key Key) (value Value, ok bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	if item, ok := cache.items[key]; ok {
		cache.list.MoveToFront(item)
		return item.Value.(*Item[Key, Value]).value, true
	}

	return
}

func (cache *LRUCache[Key, Value]) Set(key Key, value Value) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	if element, exists := cache.items[key]; exists {
		cache.list.MoveToFront(element)
		element.Value.(*Item[Key, Value]).value = value
	}

	if cache.list.Len() == cache.capacity {
		cache.remove()
	}

	item := &Item[Key, Value]{
		key:   key,
		value: value,
	}

	element := cache.list.PushFront(item)
	cache.items[key] = element
}

func (cache *LRUCache[Key, Value]) remove() {
	if element := cache.list.Back(); element != nil {
		item := cache.list.Remove(element).(*Item[Key, Value])
		delete(cache.items, item.key)
	}
}
