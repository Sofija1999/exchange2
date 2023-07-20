package usecases

import (
	"context"
	"egw-be/sofija/core/ports"
	"fmt"

	domain "github.com/Bloxico/exchange-gateway/sofija/core/domain"
	"github.com/Bloxico/exchange-gateway/sofija/repo"
	"github.com/pkg/errors"
)

// Check this service satisfies interface
var _ ports.EgwProductUsecase = (*EgwProductService)(nil)

type EgwProductService struct {
	productRepo *repo.EgwProductRepository
}

func NewEgwProductService(productRepo *repo.EgwProductRepository) *EgwProductService {
	return &EgwProductService{
		productRepo: productRepo,
	}
}

func (s *EgwProductService) InsertProduct(ctx context.Context, product *domain.EgwProduct) error {
	err := s.productRepo.Insert(ctx, product)
	if err != nil {
		return errors.Wrap(err, "Failed to insert product")
	}
	return nil
}

func (s *EgwProductService) FindByID(ctx context.Context, id string) (*domain.EgwProduct, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve product")
	}

	return product, nil
}

func (s *EgwProductService) Update(ctx context.Context, egwProduct *domain.EgwProduct) error {
	err := s.productRepo.Update(ctx, egwProduct)
	if err != nil {
		return err
	}
	return nil
}

func (s *EgwProductService) Delete(ctx context.Context, id string) error {
	err := s.productRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *EgwProductService) FindByName(ctx context.Context, name string) (*domain.EgwProduct, error) {
	product, err := s.productRepo.FindByName(ctx, name)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve product")
	}

	return product, nil
}

func (s *EgwProductService) GetAll(ctx context.Context) ([]*domain.EgwProduct, error) {
	products, err := s.productRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve products: %w", err)
	}
	return products, nil
}

func (s *EgwProductService) GetProduct(ctx context.Context, id string) (*domain.EgwProduct, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve product with ID %s: %w", id, err)
	}

	// Check if the product is nil, which indicates that the product with the given ID was not found
	if product == nil {
		return nil, nil
	}

	return product, nil
}
