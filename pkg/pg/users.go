package pg

import (
	"context"
	"cultured/pkg/models"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *pgxpool.Pool
}

func (m *UserModel) CreateUser(ctx context.Context, name, email, password string) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return 0, nil
	}

	stmt := `INSERT INTO users (name, email, password_hash)
        VALUES ($1, $2, $3)
        RETURNING id`

	var id int64
	err = m.DB.QueryRow(ctx, stmt, name, email, string(hashedPassword)).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_key" {
			return 0, models.ErrDuplicateEmail
		}
	}

	return id, nil

}

func (m *UserModel) Authenticate(ctx context.Context, email, password string) (int64, error) {
	return 0, nil
}

func (m *UserModel) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	return nil, nil
}
