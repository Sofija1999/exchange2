package repo

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/Bloxico/exchange-gateway/sofija/core/domain"
	"github.com/Bloxico/exchange-gateway/sofija/core/ports"
	"github.com/Bloxico/exchange-gateway/sofija/database"
	"github.com/google/uuid"
)

var ErrEgwProductNotFound = errors.New("product not found")

// Verify the impl matches the interface
var _ ports.EgwProductRepo = (*EgwProductRepository)(nil)

type EgwProductRepository struct {
	db *database.DB
}

func NewEgwProductRepository(db *database.DB) *EgwProductRepository {
	return &EgwProductRepository{
		db: db,
	}
}

func (repo *EgwProductRepository) Insert(ctx context.Context, EgwProduct *domain.EgwProduct) error {

	uuid := uuid.New().String()

	_, err := repo.db.Exec(ctx,
		"INSERT INTO egw.product (id, name, short_description, description, price) VALUES ($1, $2, $3, $4, $5)",
		uuid, EgwProduct.Name, EgwProduct.ShortDescription, EgwProduct.Description, EgwProduct.Price)
	if err != nil {
		return err
	}

	return nil
}

func (repo *EgwProductRepository) Update(ctx context.Context, EgwProduct *domain.EgwProduct) error {
	// update, and reflect changes in the struct
	err := repo.db.QueryRow(ctx,
		`UPDATE egw.product SET
			short_description = $1,
			description = $2,
			price = $3
		 WHERE id = $4
		 RETURNING id,name, short_description, description, price`,
		EgwProduct.ShortDescription, EgwProduct.Description, EgwProduct.Price, EgwProduct.ID).StructScan(EgwProduct)
	if err != nil {
		return err
	}

	return nil
}

func (repo *EgwProductRepository) Delete(ctx context.Context, productID string) error {
	_, err := repo.db.Exec(ctx, "DELETE FROM egw.product WHERE id = $1", productID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *EgwProductRepository) FindByID(ctx context.Context, id string) (*domain.EgwProduct, error) {
	var EgwProduct domain.EgwProduct

	err := repo.db.
		QueryRow(ctx, `SELECT id, name, short_description, description, price FROM egw.product WHERE id = $1`, id).
		StructScan(&EgwProduct)
	if err == sql.ErrNoRows {
		return nil, ErrEgwProductNotFound
	}
	if err != nil {
		return nil, err
	}

	return &EgwProduct, nil
}

func (repo *EgwProductRepository) FindByName(ctx context.Context, name string) (*domain.EgwProduct, error) {
	var EgwProduct domain.EgwProduct

	err := repo.db.QueryRow(ctx,
		`SELECT id, name, short_description, description, price FROM egw.product WHERE name = $1`,
		name).
		StructScan(&EgwProduct)
	if err == sql.ErrNoRows {
		return nil, ErrEgwProductNotFound
	}
	if err != nil {
		return nil, err
	}

	return &EgwProduct, nil
}

func (repo *EgwProductRepository) GetAll(ctx context.Context) ([]*domain.EgwProduct, error) {
	rows, err := repo.db.Query(ctx, "SELECT id, name, short_description, description, price FROM egw.product")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*domain.EgwProduct

	for rows.Next() {
		var product domain.EgwProduct
		err := rows.Scan(&product.ID, &product.Name, &product.ShortDescription, &product.Description, &product.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
