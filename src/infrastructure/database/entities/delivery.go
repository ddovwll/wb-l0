package entities

import (
	"demoService/src/domain/models"

	"github.com/jackc/pgx/v5/pgtype"
)

type Delivery struct {
	ID      pgtype.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name    string
	Phone   string
	Zip     string
	City    string
	Address string
	Region  string
	Email   string
}

func (delivery Delivery) MapToModel() *models.Delivery {
	return &models.Delivery{
		Name:    delivery.Name,
		Phone:   delivery.Phone,
		Zip:     delivery.Zip,
		City:    delivery.City,
		Address: delivery.Address,
		Region:  delivery.Region,
		Email:   delivery.Email,
	}
}

func MapDeliveryToEntity(delivery models.Delivery) *Delivery {
	return &Delivery{
		Name:    delivery.Name,
		Phone:   delivery.Phone,
		Zip:     delivery.Zip,
		City:    delivery.City,
		Address: delivery.Address,
		Region:  delivery.Region,
		Email:   delivery.Email,
	}
}
