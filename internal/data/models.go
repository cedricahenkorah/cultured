package data

import (
	"errors"
	"time"

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

func New(db *pgxpool.Pool, queryTimeout time.Duration) Models {
	return Models{
		Reviews: ReviewModel{DB: db, queryTimeout: queryTimeout},
		Users:   UserModel{DB: db, queryTimeout: queryTimeout},
	}
}
