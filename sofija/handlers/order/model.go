package order

import (
	"time"

	"github.com/Bloxico/exchange-gateway/sofija/core/domain"

	"github.com/Bloxico/exchange-gateway/sofija/core/order"
	"github.com/Bloxico/exchange-gateway/sofija/core/order_item"
)

type EgwOrderModel struct {
	ID        string                          `json:"id"`
	UserID    string                          `json:"user_id"`
	Status    string                          `json:"status"`
	CreatedAt time.Time                       `json:"created_at"`
	UpdatedAt time.Time                       `json:"updated_at"`
	Items     []*order_item.EgwItemOrderModel `json:"order_items"`
}

func (e *EgwOrderModel) FromDomain(egwOrder *domain.EgwOrder) {
	// sanity checks
	if e == nil || egwOrder == nil {
		return
	}

	e.ID = egwOrder.ID
	e.UserID = egwOrder.UserID
	e.Status = egwOrder.Status
	e.CreatedAt = egwOrder.CreatedAt
	e.UpdatedAt = egwOrder.UpdatedAt

	// Convert domain.EgwItemOrderModel to []*domain.EgwItemOrderModel
	e.Items = make([]order_item.EgwItemOrderModel, len(egwOrder.OrderItems))
	for i, item := range egwOrder.OrderItems {
		e.Items[i] = &order_item.EgwItemOrderModel{
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

	egwItemOrders := make([]*order.EgwItemOrder, len(e.Items))
	for i, item := range e.Items {
		egwItemOrders[i] = &order.EgwItemOrder{
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
