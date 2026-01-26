package service

import (
	"context"
	"net/url"
	"strings"
	"time"

	db "url_shortner_backend/db/output"
	"url_shortner_backend/internal/shorturl/repo"
	"url_shortner_backend/pkg/logger"
)

type CreateShortURLInput struct {
	LongURL   string
	OwnerType string // "user" or "anonymous"
	OwnerID   string // user_id or anon_id
}

type CreateShortURLOutput struct {
	ID        int64
	Code      string
	LongURL   string
	OwnerType string
}

type GetLongURLInput struct {
	Code string
}

type GetLongURLOutput struct {
	LongURL string
}

type GetMyURLsInput struct {
	OwnerType string
	OwnerID   string
}

type ShortURLItem struct {
	ID        int64
	Code      string
	LongURL   string
	IsActive  bool
	CreatedAt time.Time
}

type GetMyURLsOutput struct {
	URLs []ShortURLItem
}

type TransferURLsInput struct {
	AnonID string
	UserID string
}

type ShortURLSvc interface {
	CreateShortURL(ctx context.Context, input CreateShortURLInput) (CreateShortURLOutput, error)
	GetLongURL(ctx context.Context, input GetLongURLInput) (GetLongURLOutput, error)
	GetMyURLs(ctx context.Context, input GetMyURLsInput) (GetMyURLsOutput, error)
	TransferURLsToUser(ctx context.Context, input TransferURLsInput) error
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
	// Normalize URL - add https:// if no scheme provided
	longURL := normalizeURL(input.LongURL)

	// Validate URL
	if _, err := url.ParseRequestURI(longURL); err != nil {
		s.Logger.Error("invalid url provided", "url", longURL, "error", err)
		return CreateShortURLOutput{}, ErrInvalidURL
	}

	// Determine owner type
	ownerType := db.OwnerTypeEnumAnonymous
	if input.OwnerType == "user" {
		ownerType = db.OwnerTypeEnumUser
	}

	// Insert with temporary code
	tempCode := "temp"
	row, err := s.Repo.CreateShortURL(ctx, db.CreateShortURLParams{
		Code:      tempCode,
		LongUrl:   longURL,
		OwnerType: ownerType,
		OwnerID:   input.OwnerID,
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

	s.Logger.Info("short url created", "id", updatedRow.ID, "code", updatedRow.Code, "owner_type", updatedRow.OwnerType)

	return CreateShortURLOutput{
		ID:        updatedRow.ID,
		Code:      updatedRow.Code,
		LongURL:   updatedRow.LongUrl,
		OwnerType: string(updatedRow.OwnerType),
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

func (s *ShortURLSvcImp) GetMyURLs(ctx context.Context, input GetMyURLsInput) (GetMyURLsOutput, error) {
	if input.OwnerID == "" {
		s.Logger.Error("empty owner id provided")
		return GetMyURLsOutput{}, ErrInvalidOwner
	}

	ownerType := db.OwnerTypeEnumAnonymous
	if input.OwnerType == "user" {
		ownerType = db.OwnerTypeEnumUser
	}

	urls, err := s.Repo.GetShortURLsByOwner(ctx, db.GetShortURLsByOwnerParams{
		OwnerType: ownerType,
		OwnerID:   input.OwnerID,
	})
	if err != nil {
		s.Logger.Error("failed to get urls by owner", "owner_id", input.OwnerID, "error", err)
		return GetMyURLsOutput{}, ErrURLFetch
	}

	items := make([]ShortURLItem, len(urls))
	for i, u := range urls {
		items[i] = ShortURLItem{
			ID:        u.ID,
			Code:      u.Code,
			LongURL:   u.LongUrl,
			IsActive:  u.IsActive,
			CreatedAt: u.CreatedAt.Time,
		}
	}

	s.Logger.Info("urls retrieved", "owner_id", input.OwnerID, "count", len(items))

	return GetMyURLsOutput{URLs: items}, nil
}

func (s *ShortURLSvcImp) TransferURLsToUser(ctx context.Context, input TransferURLsInput) error {
	if input.AnonID == "" || input.UserID == "" {
		s.Logger.Error("empty anon_id or user_id")
		return ErrInvalidOwner
	}

	err := s.Repo.TransferAnonymousURLsToUser(ctx, db.TransferAnonymousURLsToUserParams{
		OwnerID:   input.AnonID,
		OwnerID_2: input.UserID,
	})
	if err != nil {
		s.Logger.Error("failed to transfer urls", "anon_id", input.AnonID, "user_id", input.UserID, "error", err)
		return ErrURLTransfer
	}

	s.Logger.Info("urls transferred", "anon_id", input.AnonID, "user_id", input.UserID)
	return nil
}

func normalizeURL(rawURL string) string {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		return "https://" + rawURL
	}
	return rawURL
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
