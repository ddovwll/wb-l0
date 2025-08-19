package tests

import (
	"context"
	"demoService/src/application/services"
	"demoService/src/domain/models"
	"demoService/src/infrastructure/cache"
	"demoService/src/tests/mocks"
	"testing"
)

func TestGetOrderById_FromCache(t *testing.T) {
	repository := mocks.NewMockOrderRepository()
	lruCache, _ := cache.NewLRUCache[string, *models.Order](100)
	ctx := context.Background()

	order := models.Order{OrderUID: "cached-order"}
	lruCache.Set("cached-order", &order)

	service := services.NewOrderService(repository, lruCache)

	got, err := service.GetOrderById(ctx, "cached-order")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.OrderUID != "cached-order" {
		t.Errorf("expected order UID 'cached-order', got %s", got.OrderUID)
	}
}

func TestGetOrderById_FromRepository(t *testing.T) {
	repository := mocks.NewMockOrderRepository()
	lruCache, _ := cache.NewLRUCache[string, *models.Order](100)
	ctx := context.Background()

	order := models.Order{OrderUID: "db-order"}
	_ = repository.Create(ctx, order)

	service := services.NewOrderService(repository, lruCache)

	got, err := service.GetOrderById(ctx, "db-order")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.OrderUID != "db-order" {
		t.Errorf("expected order UID 'db-order', got %s", got.OrderUID)
	}

	if _, ok := lruCache.Get("db-order"); !ok {
		t.Errorf("expected order to be cached")
	}
}

func TestCreate_Order(t *testing.T) {
	repository := mocks.NewMockOrderRepository()
	lruCache, _ := cache.NewLRUCache[string, *models.Order](100)
	service := services.NewOrderService(repository, lruCache)
	ctx := context.Background()

	order := models.Order{OrderUID: "new-order"}
	err := service.Create(ctx, order)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := repository.GetById(ctx, "new-order")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.OrderUID != "new-order" {
		t.Errorf("expected order UID 'new-order', got %s", got.OrderUID)
	}
}

func TestPreloadOrdersInCache(t *testing.T) {
	repository := mocks.NewMockOrderRepository()
	lruCache, _ := cache.NewLRUCache[string, *models.Order](100)
	ctx := context.Background()

	_ = repository.Create(ctx, models.Order{OrderUID: "1"})
	_ = repository.Create(ctx, models.Order{OrderUID: "2"})

	service := services.NewOrderService(repository, lruCache)

	err := service.PreloadOrdersInCache(ctx, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := lruCache.Get("1"); !ok {
		t.Errorf("expected order 1 to be cached")
	}

	if _, ok := lruCache.Get("2"); !ok {
		t.Errorf("expected order 2 to be cached")
	}
}
