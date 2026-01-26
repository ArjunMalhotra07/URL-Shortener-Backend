package service

import (
	"context"
	"net/url"
	"time"

	db "url_shortner_backend/db/output"
	"url_shortner_backend/internal/shorturl/repo"
	"url_shortner_backend/pkg/logger"
)

type CreateShortURLInput struct {
	LongURL string
}

type CreateShortURLOutput struct {
	ID      int64
	Code    string
	LongURL string
}

type GetLongURLInput struct {
	Code string
}

type GetLongURLOutput struct {
	LongURL string
}

type ShortURLSvc interface {
	CreateShortURL(ctx context.Context, input CreateShortURLInput) (CreateShortURLOutput, error)
	GetLongURL(ctx context.Context, input GetLongURLInput) (GetLongURLOutput, error)
}

type ShortURLSvcImp struct {
	Logger logger.Logger
	Repo   repo.ShortURLRepository
}

func NewShortURLSvcImp(Repo repo.ShortURLRepository, Logger logger.Logger) *ShortURLSvcImp {
	return &ShortURLSvcImp{
		Repo:   Repo,
		Logger: Logger,
	}
}

func (s *ShortURLSvcImp) CreateShortURL(ctx context.Context, input CreateShortURLInput) (CreateShortURLOutput, error) {
	// Validate URL
	if _, err := url.ParseRequestURI(input.LongURL); err != nil {
		s.Logger.Error("invalid url provided", "url", input.LongURL, "error", err)
		return CreateShortURLOutput{}, ErrInvalidURL
	}

	// Insert with temporary code
	tempCode := "temp"
	row, err := s.Repo.CreateShortURL(ctx, db.CreateShortURLParams{
		Code:    tempCode,
		LongUrl: input.LongURL,
	})
	if err != nil {
		s.Logger.Error("failed to create short url", "error", err)
		return CreateShortURLOutput{}, ErrURLCreation
	}

	// Generate code from ID using base62
	code := encodeBase62(row.ID)

	// Update the row with actual code
	updatedRow, err := s.Repo.UpdateShortURLCode(ctx, db.UpdateShortURLCodeParams{
		ID:   row.ID,
		Code: code,
	})
	if err != nil {
		s.Logger.Error("failed to update short url code", "id", row.ID, "error", err)
		return CreateShortURLOutput{}, ErrURLCodeUpdate
	}

	s.Logger.Info("short url created", "id", updatedRow.ID, "code", updatedRow.Code)

	return CreateShortURLOutput{
		ID:      updatedRow.ID,
		Code:    updatedRow.Code,
		LongURL: updatedRow.LongUrl,
	}, nil
}

func (s *ShortURLSvcImp) GetLongURL(ctx context.Context, input GetLongURLInput) (GetLongURLOutput, error) {
	if input.Code == "" {
		s.Logger.Error("empty code provided")
		return GetLongURLOutput{}, ErrInvalidCode
	}

	shortURL, err := s.Repo.GetShortURLByCode(ctx, input.Code)
	if err != nil {
		s.Logger.Error("failed to get short url", "code", input.Code, "error", err)
		return GetLongURLOutput{}, ErrURLNotFound
	}

	// Check if URL is active
	if !shortURL.IsActive {
		s.Logger.Info("short url is inactive", "code", input.Code)
		return GetLongURLOutput{}, ErrURLInactive
	}

	// Check if URL has expired
	if shortURL.ExpiresAt.Valid && shortURL.ExpiresAt.Time.Before(time.Now()) {
		s.Logger.Info("short url has expired", "code", input.Code)
		return GetLongURLOutput{}, ErrURLExpired
	}

	s.Logger.Info("long url retrieved", "code", input.Code)

	return GetLongURLOutput{
		LongURL: shortURL.LongUrl,
	}, nil
}

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func encodeBase62(num int64) string {
	if num == 0 {
		return string(base62Chars[0])
	}

	var encoded []byte
	for num > 0 {
		encoded = append([]byte{base62Chars[num%62]}, encoded...)
		num /= 62
	}
	return string(encoded)
}
