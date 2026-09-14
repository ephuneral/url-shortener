package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/ephuneral/url-shortener/internal/repository"
)

var (
	ErrInvalidURL  = errors.New("invalid URL format")
	ErrURLNotFound = errors.New("URL Not Found")
)

type URLService struct {
	repo    repository.URLRepository
	baseURL string
}

func NewURLService(repo repository.URLRepository, baseURL string) *URLService {
	return &URLService{
		repo:    repo,
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (s *URLService) CreateShortURL(ctx context.Context, originalURL string) (string, error) {
	if !isValidURL(originalURL) {
		return "", ErrInvalidURL
	}

	_, shortCode, err := s.repo.Save(ctx, originalURL)
	if err != nil {
		return "", fmt.Errorf("save url in repo: %w", err)
	}

	shortUrl := fmt.Sprintf("%s/%s", s.baseURL, shortCode)

	return shortUrl, nil
}

func (s *URLService) ResolveShortURL(ctx context.Context, shortCode string) (string, error) {
	urlEntity, err := s.repo.GetByShrortCode(ctx, shortCode)
	if err != nil {
		return "", fmt.Errorf("get url by short code: %w", err)
	}

	if urlEntity == nil {
		return "", ErrURLNotFound
	}

	return urlEntity.OriginalURL, nil
}

func isValidURL(rawURL string) bool {
	if rawURL == "" {
		return false
	}

	url, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	if url.Scheme == "" || url.Host == "" {
		return false
	}

	if url.Scheme != "http" && url.Scheme != "https" {
		return false
	}

	return true
}
