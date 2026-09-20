package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Review struct {
	ID        int64
	Title     string
	Content   string
	Rating    int
	CreatedAt time.Time
	UpdatedAt time.Time
}
