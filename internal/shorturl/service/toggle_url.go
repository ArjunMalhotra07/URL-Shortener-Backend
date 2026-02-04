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
			s.Logger.Info("toggle attempt on non-owned url", "code", input.Code, "owner_id", input.OwnerID)
			return ErrURLNotOwned
		}
		s.Logger.Error("failed to verify url ownership", "code", input.Code, "error", err)
		return ErrURLToggle
	}

	err = s.Repo.ToggleURLActive(ctx, db.ToggleURLActiveParams{
		Code:      input.Code,
		OwnerType: ownerType,
		OwnerID:   input.OwnerID,
	})
	if err != nil {
		s.Logger.Error("failed to toggle url", "code", input.Code, "error", err)
		return ErrURLToggle
	}

	// Invalidate cache
	s.InvalidateURLCache(ctx, input.Code)

	s.Logger.Info("url toggled", "code", input.Code, "owner_id", input.OwnerID)
	return nil
}
