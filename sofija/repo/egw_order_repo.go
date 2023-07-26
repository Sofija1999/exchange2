package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Bloxico/exchange-gateway/sofija/core/domain"

	"github.com/Bloxico/exchange-gateway/sofija/core/ports"
	"github.com/Bloxico/exchange-gateway/sofija/database"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var ErrEgwOrderNotFound = errors.New("order not found")

// Verify the impl matches the interface
var _ ports.EgwOrderRepo = (*EgwOrderRepository)(nil)

type EgwOrderRepository struct {
	db *database.DB
}

func NewEgwOrderRepository(db *database.DB) *EgwOrderRepository {
	return &EgwOrderRepository{
		db: db,
	}
}

func (repo *EgwOrderRepository) Insert(ctx context.Context, egwOrder *domain.EgwOrder) (string, error) {
	orderID := uuid.New().String()

	_, err := repo.db.Exec(ctx,
		"INSERT INTO egw.order (id, user_id, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
		orderID, egwOrder.UserID, egwOrder.Status, egwOrder.CreatedAt, egwOrder.UpdatedAt)
	if err != nil {
		return "", err
	}

	// Insert order items
	for _, item := range egwOrder.OrderItems {

		itemID := uuid.New().String()

		_, err := repo.db.Exec(ctx,
			"INSERT INTO egw.order_item (id, order_id, product_id, product_name, quantity) VALUES ($1, $2, $3, $4, $5)",
			itemID, orderID, item.ProductID, item.ProductName, item.Quantity)
		if err != nil {
			return "", err
		}
	}

	return orderID, nil
}

func (repo *EgwOrderRepository) FindByID(ctx context.Context, id string) (*domain.EgwOrder, error) {
	var egwOrder domain.EgwOrder

	// Retrieve basic order information from the main 'egw.order' table
	err := repo.db.
		QueryRow(ctx, `SELECT id, user_id, status, created_at, updated_at FROM egw.order WHERE id = $1`, id).
		StructScan(&egwOrder)
	if err == sql.ErrNoRows {
		return nil, ErrEgwOrderNotFound
	}
	if err != nil {
		return nil, err
	}

	// Retrieve order items from the 'egw.order_item' table based on the 'order_id'
	rows, err := repo.db.Query(ctx,
		`SELECT product_id, product_name, quantity FROM egw.order_item WHERE order_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.EgwOrderItem
		err := rows.Scan(&item.ProductID, &item.ProductName, &item.Quantity)
		if err != nil {
			return nil, err
		}
		egwOrder.OrderItems = append(egwOrder.OrderItems, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &egwOrder, nil
}

func (repo *EgwOrderRepository) Delete(ctx context.Context, orderID string) error {
	_, err := repo.db.Exec(ctx, "DELETE FROM egw.order_item WHERE order_id = $1", orderID)
	if err != nil {
		fmt.Println("Error while delete order items from database:", err)
		return err
	}
	_, err = repo.db.Exec(ctx, "DELETE FROM egw.order WHERE id = $1", orderID)
	if err != nil {
		fmt.Println("Error while delete order from database", err)
		return err
	}

	return nil
}

func (repo *EgwOrderRepository) Update(ctx context.Context, EgwOrder *domain.EgwOrder) error {

	_, err := repo.db.Exec(ctx,
		`UPDATE egw.order SET
			status = $1,
			updated_at = NOW()
		 WHERE id = $2`,
		EgwOrder.Status, EgwOrder.ID)
	if err != nil {
		return err
	}
	return nil
}
