package order

import "time"

type InsertRequestData struct {
	UserID string                    `json:"user_id"`
	Status string                    `json:"status"`
	Items  []*InsertOrderItemRequest `json:"order_items"`
}

type InsertOrderItemRequest struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

type InsertResponseData struct {
	ID        string                     `json:"id"`
	UserID    string                     `json:"user_id"`
	Status    string                     `json:"status"`
	CreatedAt time.Time                  `json:"created_at"`
	UpdatedAt time.Time                  `json:"updated_at"`
	Items     []*InsertOrderItemResponse `json:"order_items"`
}

type InsertOrderItemResponse struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

type UpdateRequestData struct {
	Status string `json:"status"`
}

type UpdateResponseData struct {
	ID        string                     `json:"id"`
	UserID    string                     `json:"user_id"`
	Status    string                     `json:"status"`
	CreatedAt time.Time                  `json:"created_at"`
	UpdatedAt time.Time                  `json:"updated_at"`
	Items     []*UpdateOrderItemResponse `json:"order_items"`
}

type UpdateOrderItemResponse struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
}
