package service

import (
	"context"

	"github.com/jackc/pgx/v5"

	db "url_shortner_backend/db/output"
)

type ToggleURLInput struct {
	Code      string
	OwnerType string
	OwnerID   string
}

func (s *ShortURLSvcImp) ToggleURLActive(ctx context.Context, input ToggleURLInput) error {
	if input.Code == "" {
		return ErrInvalidCode
	}
	if input.OwnerID == "" {
		return ErrInvalidOwner
	}

	ownerType := db.OwnerTypeEnumAnonymous
	if input.OwnerType == "user" {
		ownerType = db.OwnerTypeEnumUser
	}

	// Verify ownership
	_, err := s.Repo.GetURLByCodeAndOwner(ctx, db.GetURLByCodeAndOwnerParams{
		Code:      input.Code,
		OwnerType: ownerType,
		OwnerID:   input.OwnerID,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			s.Logger.Info().Str("code", input.Code).Str("owner_id", input.OwnerID).Msg("toggle attempt on non-owned url")
			return ErrURLNotOwned
		}
		s.Logger.Err(err).Str("code", input.Code).Msg("failed to verify url ownership")
		return ErrURLToggle
	}

	err = s.Repo.ToggleURLActive(ctx, db.ToggleURLActiveParams{
		Code:      input.Code,
		OwnerType: ownerType,
		OwnerID:   input.OwnerID,
	})
	if err != nil {
		s.Logger.Err(err).Str("code", input.Code).Msg("failed to toggle url")
		return ErrURLToggle
	}

	// Invalidate cache
	s.InvalidateURLCache(ctx, input.Code)

	s.Logger.Info().Str("code", input.Code).Str("owner_id", input.OwnerID).Msg("url toggled")
	return nil
}
