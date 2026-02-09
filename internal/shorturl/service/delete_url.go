package service

import (
	"context"

	"github.com/jackc/pgx/v5"

	db "url_shortner_backend/db/output"
)

type DeleteURLInput struct {
	Code      string
	OwnerType string
	OwnerID   string
}

type DeleteURLOutput struct {
	URLsCreatedThisMonth int64
	URLsLimit            int
}

func (s *ShortURLSvcImp) DeleteURL(ctx context.Context, input DeleteURLInput) (DeleteURLOutput, error) {
	if input.Code == "" {
		return DeleteURLOutput{}, ErrInvalidCode
	}
	if input.OwnerID == "" {
		return DeleteURLOutput{}, ErrInvalidOwner
	}

	ownerType := db.OwnerTypeEnumAnonymous
	quota := s.Cfg.MonthlyQuotaAnonymous
	if input.OwnerType == "user" {
		ownerType = db.OwnerTypeEnumUser
		quota = s.Cfg.MonthlyQuotaUser
	}

	// Verify ownership and get URL ID
	url, err := s.Repo.GetURLByCodeAndOwner(ctx, db.GetURLByCodeAndOwnerParams{
		Code:      input.Code,
		OwnerType: ownerType,
		OwnerID:   input.OwnerID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			s.Logger.Info().Str("code", input.Code).Str("owner_id", input.OwnerID).Msg("delete attempt on non-owned url")
			return DeleteURLOutput{}, ErrURLNotOwned
		}
		s.Logger.Err(err).Str("code", input.Code).Msg("failed to verify url ownership")
		return DeleteURLOutput{}, ErrURLDelete
	}

	// Delete analytics/clicks for this URL
	err = s.Repo.DeleteClicksByShortURLID(ctx, url.ID)
	if err != nil {
		s.Logger.Err(err).Str("code", input.Code).Msg("failed to delete clicks")
		// Continue with URL deletion even if click deletion fails
	}

	err = s.Repo.SoftDeleteURL(ctx, db.SoftDeleteURLParams{
		Code:      input.Code,
		OwnerType: ownerType,
		OwnerID:   input.OwnerID,
	})
	if err != nil {
		s.Logger.Err(err).Str("code", input.Code).Msg("failed to delete url")
		return DeleteURLOutput{}, ErrURLDelete
	}

	// Invalidate cache
	s.InvalidateURLCache(ctx, input.Code)

	// Get updated count
	monthCount, err := s.Repo.CountURLsCreatedThisMonth(ctx, input.OwnerID)
	if err != nil {
		s.Logger.Err(err).Msg("failed to count urls after delete")
		monthCount = 0
	}

	s.Logger.Info().Str("code", input.Code).Str("owner_id", input.OwnerID).Msg("url deleted")
	return DeleteURLOutput{
		URLsCreatedThisMonth: monthCount,
		URLsLimit:            quota,
	}, nil
}
