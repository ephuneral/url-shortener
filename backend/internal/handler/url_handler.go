package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/ephuneral/url-shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

type URLHandler struct {
	service *service.URLService
}

func NewURLHandler(service *service.URLService) *URLHandler {
	return &URLHandler{service: service}
}

type CreateShortURLRequest struct {
	URL string `json:"url"`
}

type CreateShortURLResponse struct {
	ShortURL string `json:"short_url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *URLHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req CreateShortURLRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "invalid JSON request")
		return
	}

	if req.URL == "" {
		h.sendError(w, http.StatusBadRequest, "url is required")
	}

	shortURL, err := h.service.CreateShortURL(ctx, req.URL)
	if err != nil {
		if errors.Is(err, service.ErrInvalidURL) {
			h.sendError(w, http.StatusBadRequest, "invalid URL format")
			return
		}

		slog.Error("failed to create short url", "error", err, "original_url", req.URL)
		h.sendError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.sendJSON(w, http.StatusCreated, CreateShortURLResponse{
		ShortURL: shortURL,
	})
}

func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	code := chi.URLParam(r, "code")
	if code == "" {
		h.sendError(w, http.StatusBadRequest, "short code is required")
		return
	}

	originalURL, err := h.service.ResolveShortURL(ctx, code)
	if err != nil {
		if errors.Is(err, service.ErrURLNotFound) {
			h.sendError(w, http.StatusNotFound, "short url not found")
			return
		}

		slog.Error("failed to resolve short url", "error", err, "code", code)
		h.sendError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

func (h *URLHandler) sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode json response", "error", err)
	}
}

func (h *URLHandler) sendError(w http.ResponseWriter, status int, message string) {
	h.sendJSON(w, status, ErrorResponse{Error: message})
}
