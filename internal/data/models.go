package data

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

type Models struct {
	Reviews ReviewModel
	Users   UserModel
}

func New(db *pgxpool.Pool) Models {
	return Models{
		Reviews: ReviewModel{DB: db},
		Users:   UserModel{DB: db},
	}
}
