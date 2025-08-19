package services

import (
	"context"
	"demoService/src/application/contracts"
	"demoService/src/domain/contracts/repositories"
	"demoService/src/domain/models"
)

type OrderService struct {
	orderRepository repositories.OrderRepository
	cache           contracts.Cache[string, *models.Order]
}

func NewOrderService(
	orderRepository repositories.OrderRepository,
	cache contracts.Cache[string, *models.Order]) *OrderService {
	return &OrderService{
		orderRepository: orderRepository,
		cache:           cache,
	}
}

func (service *OrderService) GetOrderById(context context.Context, orderId string) (*models.Order, error) {
	if cached, ok := service.cache.Get(orderId); ok {
		return cached, nil
	}

	fromDb, err := service.orderRepository.GetById(context, orderId)
	if err == nil {
		service.cache.Set(orderId, fromDb)

		return fromDb, nil
	}

	return nil, err
}

func (service *OrderService) Create(context context.Context, order models.Order) error {
	err := service.orderRepository.Create(context, order)
	if err != nil {
		return err
	}

	return nil
}

func (service *OrderService) PreloadOrdersInCache(context context.Context, count int) error {
	orders, err := service.orderRepository.GetLatestOrders(context, count)
	if err != nil {
		return err
	}

	for _, order := range *orders {
		service.cache.Set(order.OrderUID, &order)
	}

	return nil
}
