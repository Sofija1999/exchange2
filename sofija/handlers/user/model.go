package user

import (
	"time"

	"github.com/Bloxico/exchange-gateway/sofija/core/domain"
)

type EgwUserModel struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	PasswordHash string `json:"password_hash"`

	CreatedAt time.Time `json:"created_at"`
}

func (e *EgwUserModel) FromDomain(egwUser *domain.EgwUser) {
	// sanity checks
	if e == nil || egwUser == nil {
		return
	}

	e.ID = egwUser.ID
	e.Email = egwUser.Email
	e.Name = egwUser.Name
	e.Surname = egwUser.Surname
	e.CreatedAt = egwUser.CreatedAt
	// do not populate the password hash, because we do not wish to expose that when loading from the domain
}

func (e *EgwUserModel) ToDomain() *domain.EgwUser {
	if e == nil {
		return &domain.EgwUser{}
	}

	return &domain.EgwUser{
		ID:           e.ID,
		Email:        e.Email,
		CreatedAt:    e.CreatedAt,
		Name:         e.Name,
		Surname:      e.Surname,
		PasswordHash: e.PasswordHash,
	}
}
