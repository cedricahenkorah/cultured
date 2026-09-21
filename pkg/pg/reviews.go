package pg

import (
	"context"
	"cultured/pkg/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReviewModel struct {
	DB *pgxpool.Pool
}

func (m *ReviewModel) Insert(ctx context.Context, title, content string, rating int, userID int64) (int64, error) {
	stmt := `INSERT INTO reviews (user_id, title, content, rating)
	VALUES ($1, $2, $3, $4)
	RETURNING id`

	var id int64
	err := m.DB.QueryRow(ctx, stmt, userID, title, content, rating).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *ReviewModel) Get(ctx context.Context, id int64) (*models.Review, error) {
	r := &models.Review{}

	stmt := `SELECT id, title, content, rating, created_at FROM reviews WHERE id = $1`

	err := m.DB.QueryRow(ctx, stmt, id).Scan(&r.ID, &r.Title, &r.Content, &r.Rating, &r.CreatedAt)

	if err == pgx.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return r, nil
}

func (m *ReviewModel) GetAll(ctx context.Context) ([]*models.Review, error) {
	stmt := `
    SELECT id, title, content, rating, created_at
    FROM reviews
    ORDER BY created_at DESC, id DESC Limit 10
	`

	rows, err := m.DB.Query(ctx, stmt)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	reviews := []*models.Review{}

	for rows.Next() {
		r := &models.Review{}

		err = rows.Scan(&r.ID, &r.Title, &r.Content, &r.Rating, &r.CreatedAt)

		if err != nil {
			return nil, err
		}

		reviews = append(reviews, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}
