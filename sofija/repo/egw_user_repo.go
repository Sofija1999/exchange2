package repo

import (
	"context"
	"database/sql"
	"regexp"

	"github.com/pkg/errors"

	"github.com/Bloxico/exchange-gateway/sofija/core/domain"
	"github.com/Bloxico/exchange-gateway/sofija/core/ports"
	"github.com/Bloxico/exchange-gateway/sofija/database"
)

var ErrDuplicateEmail = errors.New("email already exists")
var ErrEgwUserNotFound = errors.New("user not found")

// Verify the impl matches the interface
var _ ports.EgwUserRepo = (*EgwUserRepository)(nil)

type EgwUserRepository struct {
	db *database.DB
}

func NewEgwUserRepository(db *database.DB) *EgwUserRepository {
	return &EgwUserRepository{
		db: db,
	}
}

func (repo *EgwUserRepository) Insert(ctx context.Context, EgwUser *domain.EgwUser) error {
	_, err := repo.db.Exec(ctx,
		"INSERT INTO egw.user (email, first_name, surname,password_hash) VALUES ($1, $2, $3, $4)",
		EgwUser.Email, EgwUser.Name, EgwUser.Surname, EgwUser.PasswordHash)
	if err != nil {
		alreadyExists, _ := regexp.Match(`user_email_key`, []byte(err.Error()))
		if alreadyExists {
			return ErrDuplicateEmail
		}
		return err
	}

	return nil
}

func (repo *EgwUserRepository) Update(ctx context.Context, EgwUser *domain.EgwUser) error {
	// update, and reflect changes in the struct
	err := repo.db.QueryRow(ctx,
		`UPDATE egw.user SET
			first_name = $1,
			surname = $2
		 WHERE id = $3
		 RETURNING id, first_name, surname, email`,
		EgwUser.Name, EgwUser.Surname, EgwUser.ID).StructScan(EgwUser)
	if err != nil {
		return err
	}

	return nil
}

func (repo *EgwUserRepository) FindByID(ctx context.Context, id string) (*domain.EgwUser, error) {
	var EgwUser domain.EgwUser

	err := repo.db.
		QueryRow(ctx, `SELECT id, email, first_name, surname, password_hash FROM egw.user WHERE id = $1`, id).
		StructScan(&EgwUser)
	if err == sql.ErrNoRows {
		return nil, ErrEgwUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &EgwUser, nil
}

func (repo *EgwUserRepository) FindByEmail(ctx context.Context, email string) (*domain.EgwUser, error) {
	var EgwUser domain.EgwUser

	err := repo.db.QueryRow(ctx,
		`SELECT id, email, first_name, surname, password_hash FROM egw.user WHERE email = $1`,
		email).
		StructScan(&EgwUser)
	if err == sql.ErrNoRows {
		return nil, ErrEgwUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &EgwUser, nil
}
