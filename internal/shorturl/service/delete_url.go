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

	// Verify ownership
	_, err := s.Repo.GetURLByCodeAndOwner(ctx, db.GetURLByCodeAndOwnerParams{
		Code:      input.Code,
		OwnerType: ownerType,
		OwnerID:   input.OwnerID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			s.Logger.Info("delete attempt on non-owned url", "code", input.Code, "owner_id", input.OwnerID)
			return DeleteURLOutput{}, ErrURLNotOwned
		}
		s.Logger.Error("failed to verify url ownership", "code", input.Code, "error", err)
		return DeleteURLOutput{}, ErrURLDelete
	}

	err = s.Repo.SoftDeleteURL(ctx, db.SoftDeleteURLParams{
		Code:      input.Code,
		OwnerType: ownerType,
		OwnerID:   input.OwnerID,
	})
	if err != nil {
		s.Logger.Error("failed to delete url", "code", input.Code, "error", err)
		return DeleteURLOutput{}, ErrURLDelete
	}

	// Invalidate cache
	s.InvalidateURLCache(ctx, input.Code)

	// Get updated count
	monthCount, err := s.Repo.CountURLsCreatedThisMonth(ctx, input.OwnerID)
	if err != nil {
		s.Logger.Error("failed to count urls after delete", "error", err)
		monthCount = 0
	}

	s.Logger.Info("url deleted", "code", input.Code, "owner_id", input.OwnerID)
	return DeleteURLOutput{
		URLsCreatedThisMonth: monthCount,
		URLsLimit:            quota,
	}, nil
}
