package order

import (
	"time"

	"github.com/Bloxico/exchange-gateway/sofija/core/domain"
)

type EgwOrderModel struct {
	ID        string               `json:"id"`
	UserID    string               `json:"user_id"`
	Status    string               `json:"status"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
	Items     []*EgwItemOrderModel `json:"order_items"`
}

func (e *EgwOrderModel) FromDomain(egwOrder *domain.EgwOrder) {
	if e == nil || egwOrder == nil {
		return
	}

	e.ID = egwOrder.ID
	e.UserID = egwOrder.UserID
	e.Status = egwOrder.Status
	e.CreatedAt = egwOrder.CreatedAt
	e.UpdatedAt = egwOrder.UpdatedAt

	e.Items = make([]*EgwItemOrderModel, len(egwOrder.OrderItems))
	for i, item := range egwOrder.OrderItems {
		e.Items[i] = &EgwItemOrderModel{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
		}
	}
}

func (e *EgwOrderModel) ToDomain() *domain.EgwOrder {
	if e == nil {
		return &domain.EgwOrder{}
	}

	egwItemOrders := make([]*domain.EgwOrderItem, len(e.Items))
	for i, item := range e.Items {
		egwItemOrders[i] = &domain.EgwOrderItem{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
		}
	}

	return &domain.EgwOrder{
		ID:         e.ID,
		UserID:     e.UserID,
		Status:     e.Status,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
		OrderItems: egwItemOrders,
	}
}

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
