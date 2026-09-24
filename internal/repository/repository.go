package repository

import (
	"context"

	"github.com/ephuneral/url-shortener/internal/model"
)

type URLRepository interface {
	Save(ctx context.Context, originalURL string) (id int64, shortCode string, err error)
	GetByShrortCode(ctx context.Context, shortCode string) (*model.URL, error)
}
