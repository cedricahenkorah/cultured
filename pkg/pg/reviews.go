package pg

import (
	"cultured/pkg/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReviewModel struct {
	DB *pgxpool.Pool
}

func (m *ReviewModel) Insert(title, content string, rating int) (int, error) {
	return 0, nil
}

func (m *ReviewModel) Get(id int) (*models.Review, error) {
	return nil, nil
}

func (m *ReviewModel) GetAll() ([]*models.Review, error) {
	return nil, nil
}
