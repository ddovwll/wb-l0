package tests

import (
	"demoService/src/infrastructure/cache"
	"testing"
)

func TestNewLRUCache_InvalidCapacity(t *testing.T) {
	_, err := cache.NewLRUCache[string, string](0)
	if err == nil {
		t.Errorf("expected error for capacity <= 0")
	}
}

func TestSetAndGet(t *testing.T) {
	lruCache, _ := cache.NewLRUCache[string, int](2)

	lruCache.Set("a", 1)
	lruCache.Set("b", 2)

	if value, ok := lruCache.Get("a"); value != 1 || !ok {
		t.Errorf("expected 1 and true, got %v and %v", value, ok)
	}
	if value, ok := lruCache.Get("b"); value != 2 || !ok {
		t.Errorf("expected 2 and true, got %v and %v", value, ok)
	}
}

func TestUpdateValue(t *testing.T) {
	lruCache, _ := cache.NewLRUCache[string, int](2)

	lruCache.Set("a", 1)
	lruCache.Set("a", 100)

	if value, ok := lruCache.Get("a"); value != 100 || !ok {
		t.Errorf("expected 100 and true, got %v and %v", value, ok)
	}
}

func TestEviction(t *testing.T) {
	lruCache, _ := cache.NewLRUCache[string, int](2)

	lruCache.Set("a", 1)
	lruCache.Set("b", 2)
	lruCache.Set("c", 3)

	if _, ok := lruCache.Get("a"); ok {
		t.Errorf("expected not ok, got ok")
	}
	if value, ok := lruCache.Get("b"); value != 2 || !ok {
		t.Errorf("expected 2 and true, got %v and %v", value, ok)
	}
	if value, ok := lruCache.Get("c"); value != 3 || !ok {
		t.Errorf("expected 3 and true, got %v and %v", value, ok)
	}
}

func TestLRUOrder(t *testing.T) {
	lruCache, _ := cache.NewLRUCache[string, int](2)

	lruCache.Set("a", 1)
	lruCache.Set("b", 2)
	_, _ = lruCache.Get("a")
	lruCache.Set("c", 3)

	if _, ok := lruCache.Get("b"); ok {
		t.Errorf("expected false, got %v", ok)
	}
	if value, ok := lruCache.Get("a"); value != 1 || !ok {
		t.Errorf("expected 1 and true, got %v and %v", value, ok)
	}
	if value, ok := lruCache.Get("c"); value != 3 || !ok {
		t.Errorf("expected 3 and true, got %v and %v", value, ok)
	}
}
