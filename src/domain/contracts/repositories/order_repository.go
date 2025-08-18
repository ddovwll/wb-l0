package repositories

import (
	"context"
	"demoService/src/domain/models"
)

type OrderRepository interface {
	Create(context context.Context, order models.Order) error
	GetById(context context.Context, id string) (*models.Order, error)
	GetLatestOrders(context context.Context, count int) (*[]models.Order, error)
}
