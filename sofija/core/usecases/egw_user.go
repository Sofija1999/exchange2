package usecases

import (
	"context"

	domain "github.com/Bloxico/exchange-gateway/sofija/core/domain"
	ports "github.com/Bloxico/exchange-gateway/sofija/core/ports"
	"github.com/Bloxico/exchange-gateway/sofija/repo"
	"github.com/pkg/errors"
)

// Check this service satisfies interface
var _ ports.EgwUserUsecase = (*EgwUserService)(nil)

type EgwUserService struct {
	userRepo *repo.EgwUserRepository
}

func NewEgwUserService(userRepo *repo.EgwUserRepository) *EgwUserService {
	return &EgwUserService{
		userRepo: userRepo,
	}
}

func (s *EgwUserService) RegisterUser(ctx context.Context, user *domain.EgwUser) error {
	err := s.userRepo.Insert(ctx, user)
	if err != nil {
		return errors.Wrap(err, "Failed to register user")
	}
	return nil
}

func (s *EgwUserService) FindByID(ctx context.Context, id string) (*domain.EgwUser, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve user")
	}

	return user, nil
}

func (s *EgwUserService) FindByEmail(ctx context.Context, email string) (*domain.EgwUser, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve user")
	}

	return user, nil
}

func (s *EgwUserService) Update(ctx context.Context, egwUser *domain.EgwUser) error {
	err := s.userRepo.Update(ctx, egwUser)
	if err != nil {
		return err
	}
	return nil
}
