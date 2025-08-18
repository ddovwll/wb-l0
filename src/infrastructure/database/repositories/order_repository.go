package repositories

import (
	"context"
	"demoService/src/domain"
	"demoService/src/domain/models"
	"demoService/src/infrastructure/database/entities"
	"errors"

	"gorm.io/gorm"
)

type OrderRepository struct {
	Db gorm.DB
}

func NewOrderRepository(db gorm.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (repository *OrderRepository) Create(context context.Context, order models.Order) error {
	return repository.Db.WithContext(context).Transaction(func(tx *gorm.DB) error {
		orderEntity := entities.MapOrderToEntity(order)

		if err := tx.Create(&orderEntity).Error; err != nil {
			return err
		}

		return nil
	})
}

func (repository *OrderRepository) GetById(context context.Context, id string) (*models.Order, error) {
	var order entities.Order
	result := repository.Db.
		WithContext(context).
		Preload("Delivery").
		Preload("Payment").
		Preload("Items").
		First(&order, "order_uid = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, result.Error
	}

	return order.MapToModel(), nil
}

func (repository *OrderRepository) GetLatestOrders(context context.Context, count int) (*[]models.Order, error) {
	var orders []entities.Order
	result := repository.Db.
		WithContext(context).
		Preload("Payment").
		Preload("Items").
		Preload("Delivery").
		Order("date_created desc").
		Limit(count).
		Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}

	orderModels := make([]models.Order, len(orders))
	for i, order := range orders {
		orderModels[i] = *order.MapToModel()
	}

	return &orderModels, nil
}
