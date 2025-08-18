package mocks

import (
	"context"
	"demoService/src/domain"
	"demoService/src/domain/models"
	"sync"
)

type MockOrderRepository struct {
	orders map[string]models.Order
	mutex  sync.Mutex
}

func NewMockOrderRepository() *MockOrderRepository {
	return &MockOrderRepository{
		orders: make(map[string]models.Order),
	}
}

func (orderRepository *MockOrderRepository) Create(_ context.Context, order models.Order) error {
	orderRepository.mutex.Lock()
	defer orderRepository.mutex.Unlock()

	orderRepository.orders[order.OrderUID] = order
	return nil
}

func (orderRepository *MockOrderRepository) GetById(_ context.Context, id string) (*models.Order, error) {
	orderRepository.mutex.Lock()
	defer orderRepository.mutex.Unlock()

	order, ok := orderRepository.orders[id]
	if !ok {
		return nil, domain.ErrNotFound
	}

	return &order, nil
}

func (orderRepository *MockOrderRepository) GetLatestOrders(_ context.Context, count int) (*[]models.Order, error) {
	orderRepository.mutex.Lock()
	defer orderRepository.mutex.Unlock()

	orders := make([]models.Order, 0, len(orderRepository.orders))
	for _, order := range orderRepository.orders {
		orders = append(orders, order)
	}

	if count < len(orders) {
		orders = orders[:count]
	}

	return &orders, nil
}
