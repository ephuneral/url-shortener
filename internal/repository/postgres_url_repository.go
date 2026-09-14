package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/ephuneral/url-shortener/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresURLRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresURLRepository(pool *pgxpool.Pool) *PostgresURLRepository {
	return &PostgresURLRepository{pool: pool}
}

func (r *PostgresURLRepository) Save(ctx context.Context, originalURL string) (int64, error) {
	query := `INSERT INTO urls (original_url) VALUES ($1) RETURNING id`

	var id int64
	err := r.pool.QueryRow(ctx, query, originalURL).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("save url: %w", err)
	}

	return id, nil
}

func (r *PostgresURLRepository) GetByShrortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	query := `SELECT id, short_code, original_url, created_at FROM urls WHERE short_code = $1`

	url := &model.URL{}
	err := r.pool.QueryRow(ctx, query, shortCode).Scan(
		&url.ID,
		&url.ShortCode,
		&url.OriginalURL,
		&url.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get url by short code %q: %w", shortCode, err)
	}

	return url, nil
}
