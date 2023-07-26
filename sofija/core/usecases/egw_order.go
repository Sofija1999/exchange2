package usecases

import (
	"context"
	"egw-be/sofija/core/ports"

	"github.com/Bloxico/exchange-gateway/sofija/core/domain"

	"github.com/Bloxico/exchange-gateway/sofija/repo"
	"github.com/pkg/errors"
)

// Check this service satisfies interface
var _ ports.EgwOrderUsecase = (*EgwOrderService)(nil)

type EgwOrderService struct {
	orderRepo *repo.EgwOrderRepository
}

func NewEgwOrderService(orderRepo *repo.EgwOrderRepository) *EgwOrderService {
	return &EgwOrderService{
		orderRepo: orderRepo,
	}
}

func (s *EgwOrderService) InsertOrder(ctx context.Context, order *domain.EgwOrder) (string, error) {
	insertedOrderID, err := s.orderRepo.Insert(ctx, order)
	if err != nil {
		return "", errors.Wrap(err, "Failed to insert order")
	}

	return insertedOrderID, nil
}

func (s *EgwOrderService) FindByID(ctx context.Context, id string) (*domain.EgwOrder, error) {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve order")
	}

	return order, nil
}

func (s *EgwOrderService) Delete(ctx context.Context, id string) error {
	err := s.orderRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *EgwOrderService) Update(ctx context.Context, egwOrder *domain.EgwOrder) error {
	err := s.orderRepo.Update(ctx, egwOrder)
	if err != nil {
		return err
	}
	return nil
}
