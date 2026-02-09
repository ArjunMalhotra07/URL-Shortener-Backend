package service

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	db "url_shortner_backend/db/output"
)

type UpdateLongURLInput struct {
	Code      string
	LongURL   string
	ExpiresAt *time.Time // nil means no change, zero time means remove expiry
	OwnerType string
	OwnerID   string
}

type UpdateLongURLOutput struct {
	Code      string
	LongURL   string
	IsActive  bool
	ExpiresAt *time.Time
	IsExpired bool
	CreatedAt time.Time
}

func (s *ShortURLSvcImp) UpdateLongURL(ctx context.Context, input UpdateLongURLInput) (*UpdateLongURLOutput, error) {
	if input.Code == "" {
		return nil, ErrInvalidCode
	}
	if input.OwnerID == "" {
		return nil, ErrInvalidOwner
	}

	longURL := normalizeURL(input.LongURL)
	if err := validateURL(longURL); err != nil {
		s.Logger.Err(err).Str("url", longURL).Msg("invalid url provided")
		return nil, ErrInvalidURL
	}

	// Validate expires_at if provided
	if input.ExpiresAt != nil && !input.ExpiresAt.IsZero() && input.ExpiresAt.Before(time.Now()) {
		return nil, ErrInvalidExpiresAt
	}

	ownerType := db.OwnerTypeEnumAnonymous
	if input.OwnerType == "user" {
		ownerType = db.OwnerTypeEnumUser
	}

	// Verify ownership and get current URL state
	existingURL, err := s.Repo.GetURLByCodeAndOwner(ctx, db.GetURLByCodeAndOwnerParams{
		Code:      input.Code,
		OwnerType: ownerType,
		OwnerID:   input.OwnerID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			s.Logger.Info().Str("owner_id", input.OwnerID).Str("code", input.Code).Msg("update attempt on non-owned url")
			return nil, ErrURLNotOwned
		}
		s.Logger.Err(err).Str("code", input.Code).Msg("failed to verify url ownership")
		return nil, ErrURLUpdate
	}

	// Determine expires_at value for update
	var expiresAt pgtype.Timestamptz
	if input.ExpiresAt != nil {
		if input.ExpiresAt.IsZero() {
			// Zero time means remove expiry
			expiresAt = pgtype.Timestamptz{Valid: false}
		} else {
			expiresAt = pgtype.Timestamptz{Time: *input.ExpiresAt, Valid: true}
		}
	} else {
		// No change - keep existing value
		expiresAt = existingURL.ExpiresAt
	}

	updatedURL, err := s.Repo.UpdateLongURL(ctx, db.UpdateLongURLParams{Code: input.Code, OwnerType: ownerType, OwnerID: input.OwnerID, LongUrl: longURL, ExpiresAt: expiresAt})
	if err != nil {
		s.Logger.Err(err).Str("code", input.Code).Msg("failed to update long url")
		return nil, ErrURLUpdate
	}
	// Invalidate cache
	s.InvalidateURLCache(ctx, input.Code)
	s.Logger.Info().Str("owner_id", input.OwnerID).Str("code", input.Code).Msg("url updated")

	// Build response
	output := &UpdateLongURLOutput{
		Code:      updatedURL.Code,
		LongURL:   updatedURL.LongUrl,
		IsActive:  updatedURL.IsActive,
		CreatedAt: updatedURL.CreatedAt.Time,
	}

	if updatedURL.ExpiresAt.Valid {
		output.ExpiresAt = &updatedURL.ExpiresAt.Time
		output.IsExpired = updatedURL.ExpiresAt.Time.Before(time.Now())
	}

	return output, nil
}
