package tests

import (
	"demoService/src/infrastructure/cache"
	"testing"
)

func TestNewLRUCache_InvalidCapacity(t *testing.T) {
	_, err := cache.NewLRUCache(0)
	if err == nil {
		t.Errorf("expected error for capacity <= 0")
	}
}

func TestSetAndGet(t *testing.T) {
	lruCache, _ := cache.NewLRUCache(2)

	lruCache.Set("a", 1)
	lruCache.Set("b", 2)

	if v := lruCache.Get("a"); v != 1 {
		t.Errorf("expected 1, got %v", v)
	}
	if v := lruCache.Get("b"); v != 2 {
		t.Errorf("expected 2, got %v", v)
	}
}

func TestUpdateValue(t *testing.T) {
	lruCache, _ := cache.NewLRUCache(2)

	lruCache.Set("a", 1)
	lruCache.Set("a", 100)

	if v := lruCache.Get("a"); v != 100 {
		t.Errorf("expected 100, got %v", v)
	}
}

func TestEviction(t *testing.T) {
	lruCache, _ := cache.NewLRUCache(2)

	lruCache.Set("a", 1)
	lruCache.Set("b", 2)
	lruCache.Set("c", 3)

	if v := lruCache.Get("a"); v != nil {
		t.Errorf("expected nil, got %v", v)
	}
	if v := lruCache.Get("b"); v != 2 {
		t.Errorf("expected 2, got %v", v)
	}
	if v := lruCache.Get("c"); v != 3 {
		t.Errorf("expected 3, got %v", v)
	}
}

func TestLRUOrder(t *testing.T) {
	lruCache, _ := cache.NewLRUCache(2)

	lruCache.Set("a", 1)
	lruCache.Set("b", 2)
	_ = lruCache.Get("a")
	lruCache.Set("c", 3)

	if v := lruCache.Get("b"); v != nil {
		t.Errorf("expected nil (b should be evicted), got %v", v)
	}
	if v := lruCache.Get("a"); v != 1 {
		t.Errorf("expected 1, got %v", v)
	}
	if v := lruCache.Get("c"); v != 3 {
		t.Errorf("expected 3, got %v", v)
	}
}
