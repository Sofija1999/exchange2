package ports

import (
	"context"

	domain "github.com/Bloxico/exchange-gateway/sofija/core/domain"
)

type EgwUserUsecase interface {
	RegisterUser(ctx context.Context, user *domain.EgwUser) error
	FindByID(ctx context.Context, id string) (*domain.EgwUser, error)
	FindByEmail(ctx context.Context, email string) (*domain.EgwUser, error)
	Update(ctx context.Context, egwUser *domain.EgwUser) error
}

type EgwProductUsecase interface {
	InsertProduct(ctx context.Context, product *domain.EgwProduct) error
	FindByID(ctx context.Context, id string) (*domain.EgwProduct, error)
	Update(ctx context.Context, EgwProduct *domain.EgwProduct) error
	Delete(ctx context.Context, id string) error
	FindByName(ctx context.Context, name string) (*domain.EgwProduct, error)
	GetAll(ctx context.Context) ([]*domain.EgwProduct, error)
	GetProduct(ctx context.Context, id string) (*domain.EgwProduct, error)
}

type EgwOrderUsecase interface {
	InsertOrder(ctx context.Context, order *domain.EgwOrder) (string, error)
	FindByID(ctx context.Context, id string) (*domain.EgwOrder, error)
}
