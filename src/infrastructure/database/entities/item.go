package entities

import (
	"demoService/src/domain/models"

	"github.com/jackc/pgx/v5/pgtype"
)

type Item struct {
	ID          pgtype.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID     string      `gorm:"index"`
	ChrtID      int
	TrackNumber string
	Price       int
	RID         string
	Name        string
	Sale        int
	Size        string
	TotalPrice  int
	NmID        int
	Brand       string
	Status      int
}

func (item Item) MapToModel() *models.Item {
	return &models.Item{
		ChrtID:      item.ChrtID,
		TrackNumber: item.TrackNumber,
		Price:       item.Price,
		Rid:         item.RID,
		Name:        item.Name,
		Sale:        item.Sale,
		Size:        item.Size,
		TotalPrice:  item.TotalPrice,
		NmID:        item.NmID,
		Brand:       item.Brand,
		Status:      item.Status,
	}
}

func MapItemToEntity(item models.Item) *Item {
	return &Item{
		ChrtID:      item.ChrtID,
		TrackNumber: item.TrackNumber,
		Price:       item.Price,
		RID:         item.Rid,
		Name:        item.Name,
		Sale:        item.Sale,
		Size:        item.Size,
		TotalPrice:  item.TotalPrice,
		NmID:        item.NmID,
		Brand:       item.Brand,
		Status:      item.Status,
	}
}
