package data

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Review struct {
	ID        int64     `json:"id"`
	UserID    *int64    `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Rating    int       `json:"rating"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ReviewModel struct {
	DB *pgxpool.Pool
}

func (m *ReviewModel) Insert(ctx context.Context, title, content string, rating int, userID int64) (*Review, error) {
	query := `INSERT INTO reviews (user_id, title, content, rating)
	VALUES ($1, $2, $3, $4)
	RETURNING id, user_id, title, content, rating, created_at, updated_at`

	args := []any{userID, title, content, rating}

	r := &Review{}

	err := m.DB.QueryRow(ctx, query, args...).Scan(&r.ID, &r.UserID, &r.Title, &r.Content, &r.Rating, &r.CreatedAt, &r.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return r, nil
}

func (m *ReviewModel) Get(ctx context.Context, id int64) (*Review, error) {
	r := &Review{}

	query := `SELECT id, user_id, title, content, rating, created_at FROM reviews WHERE id = $1`

	err := m.DB.QueryRow(ctx, query, id).Scan(&r.ID, &r.UserID, &r.Title, &r.Content, &r.Rating, &r.CreatedAt)

	if err == pgx.ErrNoRows {
		return nil, ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return r, nil
}

func (m *ReviewModel) GetAll(ctx context.Context) ([]*Review, error) {
	query := `
    SELECT id, user_id, title, content, rating, created_at
    FROM reviews
    ORDER BY created_at DESC, id DESC Limit 10
	`

	rows, err := m.DB.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	reviews := []*Review{}

	for rows.Next() {
		r := &Review{}

		err = rows.Scan(&r.ID, &r.UserID, &r.Title, &r.Content, &r.Rating, &r.CreatedAt)

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
