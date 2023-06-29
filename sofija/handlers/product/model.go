package product

import (
	"time"

	"github.com/Bloxico/exchange-gateway/sofija/core/domain"
)

type EgwProductModel struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	Description      string `json:"description"`
	Price            int64  `json:"price"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (e *EgwProductModel) FromDomain(egwProduct *domain.EgwProduct) {
	// sanity checks
	if e == nil || egwProduct == nil {
		return
	}

	e.ID = egwProduct.ID
	e.Name = egwProduct.Name
	e.ShortDescription = egwProduct.ShortDescription
	e.Description = egwProduct.Description
	e.Price = egwProduct.Price
	e.CreatedAt = egwProduct.CreatedAt
	e.UpdatedAt = egwProduct.UpdatedAt
}

func (e *EgwProductModel) ToDomain() *domain.EgwProduct {
	if e == nil {
		return &domain.EgwProduct{}
	}

	return &domain.EgwProduct{
		ID:               e.ID,
		Name:             e.Name,
		ShortDescription: e.ShortDescription,
		Description:      e.Description,
		Price:            e.Price,
		CreatedAt:        e.CreatedAt,
		UpdatedAt:        e.UpdatedAt,
	}
}
