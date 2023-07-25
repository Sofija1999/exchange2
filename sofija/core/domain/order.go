package domain

import (
	"fmt"
	"time"
)

type EgwOrder struct {
	ID         string          `json:"id" db:"id"`
	UserID     string          `json:"user_id" db:"user_id"`
	Status     string          `json:"status" db:"status"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at" db:"updated_at"`
	OrderItems []*EgwOrderItem `json:"order_items" db:"-"`
}

func NewOrder(id, userID, status string) *EgwOrder {
	return &EgwOrder{
		ID:         id,
		UserID:     userID,
		Status:     status,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		OrderItems: []*EgwOrderItem{},
	}
}

func (o *EgwOrder) ToString() string {
	return fmt.Sprintf("#%s, UserID: %s, Status: %s, CreatedAt: %s, UpdatedAt: %s", o.ID, o.UserID, o.Status, o.CreatedAt, o.UpdatedAt)
}
