package ports

import (
	"context"

	domain "github.com/Bloxico/exchange-gateway/sofija/core/domain"
)

type EgwUserRepo interface {
	Insert(ctx context.Context, user *domain.EgwUser) error
	Update(ctx context.Context, user *domain.EgwUser) error
	FindByID(ctx context.Context, id string) (*domain.EgwUser, error)
	FindByEmail(ctx context.Context, email string) (*domain.EgwUser, error)
}

type EgwProductRepo interface {
	Insert(ctx context.Context, product *domain.EgwProduct) error
	Update(ctx context.Context, product *domain.EgwProduct) error
	FindByID(ctx context.Context, id string) (*domain.EgwProduct, error)
	Delete(ctx context.Context, id string) error
	FindByName(ctx context.Context, name string) (*domain.EgwProduct, error)
	GetAll(ctx context.Context) ([]*domain.EgwProduct, error)
	GetProduct(ctx context.Context, id string) (*domain.EgwProduct, error)
}

type EgwOrderRepo interface {
	Insert(ctx context.Context, order *domain.EgwOrder) (string, error)
	FindByID(ctx context.Context, id string) (*domain.EgwOrder, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, order *domain.EgwOrder) error
}
