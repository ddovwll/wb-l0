package entities

import (
	"demoService/src/domain/models"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Order struct {
	OrderUID          string `gorm:"primaryKey"`
	TrackNumber       string
	Entry             string
	DeliveryID        pgtype.UUID
	Delivery          Delivery `gorm:"foreignKey:DeliveryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PaymentID         pgtype.UUID
	Payment           Payment `gorm:"foreignKey:PaymentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Items             []Item  `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Locale            string
	InternalSignature string
	CustomerID        string
	DeliveryService   string
	ShardKey          string
	SmID              int
	DateCreated       time.Time
	OofShard          string
}

func (order Order) MapToModel() *models.Order {
	items := make([]models.Item, len(order.Items))
	for i := range items {
		items[i] = *order.Items[i].MapToModel()
	}

	return &models.Order{
		OrderUID:          order.OrderUID,
		TrackNumber:       order.TrackNumber,
		Entry:             order.Entry,
		Delivery:          *order.Delivery.MapToModel(),
		Payment:           *order.Payment.MapToModel(),
		Items:             items,
		Locale:            order.Locale,
		InternalSignature: order.InternalSignature,
		CustomerID:        order.CustomerID,
		DeliveryService:   order.DeliveryService,
		ShardKey:          order.ShardKey,
		SmID:              order.SmID,
		DateCreated:       order.DateCreated,
		OofShard:          order.OofShard,
	}
}

func MapOrderToEntity(order models.Order) *Order {
	items := make([]Item, len(order.Items))
	for i := range items {
		items[i] = *MapItemToEntity(order.Items[i])
	}

	return &Order{
		OrderUID:          order.OrderUID,
		TrackNumber:       order.TrackNumber,
		Entry:             order.Entry,
		Delivery:          *MapDeliveryToEntity(order.Delivery),
		Payment:           *MapPaymentToEntity(order.Payment),
		Items:             items,
		Locale:            order.Locale,
		InternalSignature: order.InternalSignature,
		CustomerID:        order.CustomerID,
		DeliveryService:   order.DeliveryService,
		ShardKey:          order.ShardKey,
		SmID:              order.SmID,
		DateCreated:       order.DateCreated,
		OofShard:          order.OofShard,
	}
}
