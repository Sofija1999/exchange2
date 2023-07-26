package order

type InsertRequestData struct {
	UserID string                    `json:"user_id"`
	Status string                    `json:"status"`
	Items  []*InsertOrderItemRequest `json:"order_items"`
}

type InsertOrderItemRequest struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
}
