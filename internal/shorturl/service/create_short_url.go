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
	Name      *string
}
type CreateShortURLOutput struct {
	ID                   int64
	Code                 string
	LongURL              string
	OwnerType            string
	IsActive             bool
	ExpiresAt            *time.Time
	CreatedAt            time.Time
	URLsCreatedThisMonth int64
	URLsLimit            int
	Name                 *string
}

func (s *ShortURLSvcImp) CreateShortURL(ctx context.Context, input CreateShortURLInput) (CreateShortURLOutput, error) {
	longURL := normalizeURL(input.LongURL)

	if err := validateURL(longURL); err != nil {
		s.Logger.Err(err).Str("url", longURL).Msg("invalid url provided")
		return CreateShortURLOutput{}, ErrInvalidURL
	}

	// Check monthly quota
	monthCount, err := s.Repo.CountURLsCreatedThisMonth(ctx, input.OwnerID)
	if err != nil {
		s.Logger.Err(err).Str("owner_id", input.OwnerID).Msg("failed to count this month's urls")
		return CreateShortURLOutput{}, ErrURLCreation
	}

	// Get quota based on owner type
	quota := s.Cfg.MonthlyQuotaAnonymous
	if input.OwnerType == "user" {
		quota = s.Cfg.MonthlyQuotaUser
	}

	if int(monthCount) >= quota {
		s.Logger.Info().Str("owner_id", input.OwnerID).Str("owner_type", input.OwnerType).Int64("count", monthCount).Int("quota", quota).Msg("monthly quota exceeded")
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

	var name pgtype.Text
	if input.Name != nil {
		name = pgtype.Text{String: *input.Name, Valid: true}
	}

	tempCode := "temp"
	row, err := s.Repo.CreateShortURL(ctx, db.CreateShortURLParams{
		Code:      tempCode,
		LongUrl:   longURL,
		OwnerType: ownerType,
		OwnerID:   input.OwnerID,
		ExpiresAt: expiresAt,
		Name:      name,
	})
	if err != nil {
		s.Logger.Err(err).Msg("failed to create short url")
		return CreateShortURLOutput{}, ErrURLCreation
	}

	code := encodeBase62(row.ID)

	updatedRow, err := s.Repo.UpdateShortURLCode(ctx, db.UpdateShortURLCodeParams{
		ID:   row.ID,
		Code: code,
	})
	if err != nil {
		s.Logger.Err(err).Int64("id", row.ID).Msg("failed to update short url code")
		return CreateShortURLOutput{}, ErrURLCodeUpdate
	}

	s.Logger.Info().Int64("id", updatedRow.ID).Str("code", updatedRow.Code).Str("owner_type", string(updatedRow.OwnerType)).Msg("short url created")

	var outputExpiresAt *time.Time
	if updatedRow.ExpiresAt.Valid {
		outputExpiresAt = &updatedRow.ExpiresAt.Time
	}

	var outputName *string
	if updatedRow.Name.Valid {
		outputName = &updatedRow.Name.String
	}

	return CreateShortURLOutput{
		ID:                   updatedRow.ID,
		Code:                 updatedRow.Code,
		LongURL:              updatedRow.LongUrl,
		OwnerType:            string(updatedRow.OwnerType),
		IsActive:             true,
		ExpiresAt:            outputExpiresAt,
		CreatedAt:            updatedRow.CreatedAt.Time,
		URLsCreatedThisMonth: monthCount + 1, // +1 for the one we just created
		URLsLimit:            quota,
		Name:                 outputName,
	}, nil
}
