package domain

import (
	"fmt"
	"time"
)

type EgwProduct struct {
	ID               string    `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	ShortDescription string    `json:"short_description" db:"short_description"`
	Description      string    `json:"description" db:"description"`
	Price            int64     `json:"price" db:"price"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

func NewProduct(id string, name string, short_description string, description string, price int64) *EgwProduct {
	return &EgwProduct{
		ID:               id,
		Name:             name,
		ShortDescription: short_description,
		Description:      description,
		Price:            price,
	}
}

func (e *EgwProduct) ToString() string {
	return fmt.Sprintf("#%s %s %s %s - %d", e.ID, e.Name, e.ShortDescription, e.Description, e.Price)
}
