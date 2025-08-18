package tests

import (
	"context"
	"demoService/src/application/services"
	"demoService/src/domain/models"
	"demoService/src/infrastructure/cache"
	"demoService/src/tests/mocks"
	"demoService/src/web/controllers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetOrderById_BadRequest(t *testing.T) {
	repository := mocks.NewMockOrderRepository()
	lruCache, _ := cache.NewLRUCache(100)
	service := services.NewOrderService(repository, lruCache)
	controller := controllers.NewOrderController(*service)

	req := httptest.NewRequest(http.MethodGet, "/order/", nil)
	w := httptest.NewRecorder()

	controller.GetOrderById(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetOrderById_Success(t *testing.T) {
	repository := mocks.NewMockOrderRepository()
	lruCache, _ := cache.NewLRUCache(100)
	service := services.NewOrderService(repository, lruCache)
	controller := controllers.NewOrderController(*service)

	order := models.Order{OrderUID: "123"}
	ctx := context.Background()
	_ = service.Create(ctx, order)

	req := httptest.NewRequest(http.MethodGet, "/order/123", nil)
	w := httptest.NewRecorder()

	controller.GetOrderById(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var got models.Order
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if got.OrderUID != "123" {
		t.Errorf("expected order UID 123, got %s", got.OrderUID)
	}
}

func TestGetOrderById_NotFound(t *testing.T) {
	repository := mocks.NewMockOrderRepository()
	lruCache, _ := cache.NewLRUCache(100)
	service := services.NewOrderService(repository, lruCache)
	controller := controllers.NewOrderController(*service)

	req := httptest.NewRequest(http.MethodGet, "/order/unknown", nil)
	w := httptest.NewRecorder()

	controller.GetOrderById(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}
