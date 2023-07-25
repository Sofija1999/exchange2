package domain

import (
	"fmt"
)

type EgwOrderItem struct {
	ID          string `json:"id" db:"id"`
	OrderID     string `json:"order_id" db:"order_id"`
	ProductID   string `json:"product_id" db:"product_id"`
	ProductName string `json:"product_name" db:"product_name"`
	Quantity    int    `json:"quantity" db:"quantity"`
}

func NewOrderItem(id, orderID, productID, productName string, quantity int) *EgwOrderItem {
	return &EgwOrderItem{
		ID:          id,
		OrderID:     orderID,
		ProductID:   productID,
		ProductName: productName,
		Quantity:    quantity,
	}
}

func (oi *EgwOrderItem) ToString() string {
	return fmt.Sprintf("#%s, OrderID: %s, ProductID: %s, ProductName: %s, Quantity: %d", oi.ID, oi.OrderID, oi.ProductID, oi.ProductName, oi.Quantity)
}
