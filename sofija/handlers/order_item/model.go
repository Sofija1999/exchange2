package order_item

import (
	"github.com/Bloxico/exchange-gateway/sofija/core/domain"
)

type EgwItemOrderModel struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
}

func (e *EgwItemOrderModel) FromDomain(item *domain.EgwOrderItem) {
	if e == nil || item == nil {
		return
	}

	e.ProductID = item.ProductID
	e.ProductName = item.ProductName
	e.Quantity = item.Quantity
}

func (e *EgwItemOrderModel) ToDomain() *domain.EgwOrderItem {
	if e == nil {
		return &domain.EgwOrderItem{}
	}

	return &domain.EgwOrderItem{
		ProductID:   e.ProductID,
		ProductName: e.ProductName,
		Quantity:    e.Quantity,
	}
}
