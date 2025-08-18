package entities

import (
	"demoService/src/domain/models"

	"github.com/jackc/pgx/v5/pgtype"
)

type Payment struct {
	ID           pgtype.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Transaction  string
	RequestID    string
	Currency     string
	Provider     string
	Amount       int
	PaymentDT    int64
	Bank         string
	DeliveryCost int
	GoodsTotal   int
	CustomFee    int
}

func (payment Payment) MapToModel() *models.Payment {
	return &models.Payment{
		Transaction:  payment.Transaction,
		RequestID:    payment.RequestID,
		Currency:     payment.Currency,
		Provider:     payment.Provider,
		Amount:       payment.Amount,
		PaymentDT:    payment.PaymentDT,
		Bank:         payment.Bank,
		DeliveryCost: payment.DeliveryCost,
		GoodsTotal:   payment.GoodsTotal,
		CustomFee:    payment.CustomFee,
	}
}

func MapPaymentToEntity(payment models.Payment) *Payment {
	return &Payment{
		Transaction:  payment.Transaction,
		RequestID:    payment.RequestID,
		Currency:     payment.Currency,
		Provider:     payment.Provider,
		Amount:       payment.Amount,
		PaymentDT:    payment.PaymentDT,
		Bank:         payment.Bank,
		DeliveryCost: payment.DeliveryCost,
		GoodsTotal:   payment.GoodsTotal,
		CustomFee:    payment.CustomFee,
	}
}
