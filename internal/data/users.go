package data

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserModel struct {
	DB *pgxpool.Pool
}

func (m *UserModel) CreateUser(ctx context.Context, name, email, password string) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return 0, err
	}

	stmt := `INSERT INTO users (name, email, password_hash)
        VALUES ($1, $2, $3)
        RETURNING id`

	args := []any{name, email, string(hashedPassword)}

	var id int64
	err = m.DB.QueryRow(ctx, stmt, args...).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_key" {
			return 0, ErrDuplicateEmail
		}

		return 0, err
	}

	return id, nil
}

func (m *UserModel) Authenticate(ctx context.Context, email, password string) (int64, error) {
	var (
		id           int64
		passwordHash string
	)

	stmt := `SELECT id, password_hash
        FROM users
        WHERE email = $1`

	err := m.DB.QueryRow(ctx, stmt, email).Scan(&id, &passwordHash)

	if err == pgx.ErrNoRows {
		return 0, ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))

	if err != nil {
		return 0, ErrInvalidCredentials
	}

	return id, nil
}

func (m *UserModel) GetUserByID(ctx context.Context, id int64) (*User, error) {
	u := &User{}

	stmt := `SELECT id, name, email, created_at, updated_at
    FROM users
    WHERE id = $1`

	err := m.DB.QueryRow(ctx, stmt, id).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return u, nil
}
