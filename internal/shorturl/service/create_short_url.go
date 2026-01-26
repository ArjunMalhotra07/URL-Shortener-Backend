package service

import (
	"context"
	"time"
	db "url_shortner_backend/db/output"
)

const (
	DailyQuota = 10
)

type CreateShortURLInput struct {
	LongURL   string
	OwnerType string
	OwnerID   string
}
type CreateShortURLOutput struct {
	ID        int64
	Code      string
	LongURL   string
	OwnerType string
	IsActive  bool
	CreatedAt time.Time
}

func (s *ShortURLSvcImp) CreateShortURL(ctx context.Context, input CreateShortURLInput) (CreateShortURLOutput, error) {
	longURL := normalizeURL(input.LongURL)

	if err := validateURL(longURL); err != nil {
		s.Logger.Error("invalid url provided", "url", longURL, "error", err)
		return CreateShortURLOutput{}, ErrInvalidURL
	}

	// Check daily quota
	todayCount, err := s.Repo.CountURLsCreatedToday(ctx, input.OwnerID)
	if err != nil {
		s.Logger.Error("failed to count today's urls", "owner_id", input.OwnerID, "error", err)
		return CreateShortURLOutput{}, ErrURLCreation
	}
	if todayCount >= DailyQuota {
		s.Logger.Info("daily quota exceeded", "owner_id", input.OwnerID, "count", todayCount)
		return CreateShortURLOutput{}, ErrDailyQuotaExceeded
	}

	ownerType := db.OwnerTypeEnumAnonymous
	if input.OwnerType == "user" {
		ownerType = db.OwnerTypeEnumUser
	}

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

	code := encodeBase62(row.ID)

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
		IsActive:  true,
		CreatedAt: updatedRow.CreatedAt.Time,
	}, nil
}
