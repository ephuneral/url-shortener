package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/ephuneral/url-shortener/internal/model"
	"github.com/ephuneral/url-shortener/pkg/base62"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresURLRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresURLRepository(pool *pgxpool.Pool) *PostgresURLRepository {
	return &PostgresURLRepository{pool: pool}
}

func (r *PostgresURLRepository) Save(ctx context.Context, originalURL string) (int64, string, error) {
	var id int64
	err := r.pool.QueryRow(ctx, `SELECT nextval(pg_get_serial_sequence('urls', 'id'))`).Scan(&id)
	if err != nil {
		return 0, "", fmt.Errorf("get next id: %w", err)
	}

	shortCode, err := base62.Encode(id)
	if err != nil {
		return 0, "", fmt.Errorf("encode id to base62: %w", err)
	}

	_, err = r.pool.Exec(ctx, `INSERT INTO urls (id, short_code, original_url) VALUES ($1, $2, $3)`, id, shortCode, originalURL)
	if err != nil {
		return 0, "", fmt.Errorf("insert url: %w", err)
	}

	return id, shortCode, nil
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
