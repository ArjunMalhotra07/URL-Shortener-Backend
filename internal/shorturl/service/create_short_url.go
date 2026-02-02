package service

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	db "url_shortner_backend/db/output"
)

type CreateShortURLInput struct {
	LongURL   string
	OwnerType string
	OwnerID   string
	ExpiresAt *time.Time
}
type CreateShortURLOutput struct {
	ID        int64
	Code      string
	LongURL   string
	OwnerType string
	IsActive  bool
	ExpiresAt *time.Time
	CreatedAt time.Time
}

func (s *ShortURLSvcImp) CreateShortURL(ctx context.Context, input CreateShortURLInput) (CreateShortURLOutput, error) {
	longURL := normalizeURL(input.LongURL)

	if err := validateURL(longURL); err != nil {
		s.Logger.Error("invalid url provided", "url", longURL, "error", err)
		return CreateShortURLOutput{}, ErrInvalidURL
	}

	// Check monthly quota
	monthCount, err := s.Repo.CountURLsCreatedThisMonth(ctx, input.OwnerID)
	if err != nil {
		s.Logger.Error("failed to count this month's urls", "owner_id", input.OwnerID, "error", err)
		return CreateShortURLOutput{}, ErrURLCreation
	}

	// Get quota based on owner type
	quota := s.Cfg.MonthlyQuotaAnonymous
	if input.OwnerType == "user" {
		quota = s.Cfg.MonthlyQuotaUser
	}

	if int(monthCount) >= quota {
		s.Logger.Info("monthly quota exceeded", "owner_id", input.OwnerID, "owner_type", input.OwnerType, "count", monthCount, "quota", quota)
		return CreateShortURLOutput{}, ErrMonthlyQuotaExceeded
	}

	ownerType := db.OwnerTypeEnumAnonymous
	if input.OwnerType == "user" {
		ownerType = db.OwnerTypeEnumUser
	}

	var expiresAt pgtype.Timestamptz
	if input.OwnerType == "anonymous" {
		expiresAt = pgtype.Timestamptz{Time: time.Now().Add(24 * time.Hour), Valid: true}
	} else if input.ExpiresAt != nil {
		expiresAt = pgtype.Timestamptz{Time: *input.ExpiresAt, Valid: true}
	}

	tempCode := "temp"
	row, err := s.Repo.CreateShortURL(ctx, db.CreateShortURLParams{
		Code:      tempCode,
		LongUrl:   longURL,
		OwnerType: ownerType,
		OwnerID:   input.OwnerID,
		ExpiresAt: expiresAt,
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

	var outputExpiresAt *time.Time
	if updatedRow.ExpiresAt.Valid {
		outputExpiresAt = &updatedRow.ExpiresAt.Time
	}

	return CreateShortURLOutput{
		ID:        updatedRow.ID,
		Code:      updatedRow.Code,
		LongURL:   updatedRow.LongUrl,
		OwnerType: string(updatedRow.OwnerType),
		IsActive:  true,
		ExpiresAt: outputExpiresAt,
		CreatedAt: updatedRow.CreatedAt.Time,
	}, nil
}
