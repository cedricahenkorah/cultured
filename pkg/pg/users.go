package pg

import (
	"context"
	"cultured/pkg/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserModel struct {
	DB *pgxpool.Pool
}

func (m *UserModel) CreateUser(ctx context.Context, name, email, password string) (int64, error) {
	return 0, nil
}

func (m *UserModel) Authenticate(ctx context.Context, email, password string) (int64, error) {
	return 0, nil
}

func (m *UserModel) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	return nil, nil
}
