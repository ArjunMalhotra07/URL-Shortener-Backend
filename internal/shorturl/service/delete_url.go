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

func (s *ShortURLSvcImp) DeleteURL(ctx context.Context, input DeleteURLInput) error {
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
			s.Logger.Info("delete attempt on non-owned url", "code", input.Code, "owner_id", input.OwnerID)
			return ErrURLNotOwned
		}
		s.Logger.Error("failed to verify url ownership", "code", input.Code, "error", err)
		return ErrURLDelete
	}

	err = s.Repo.SoftDeleteURL(ctx, db.SoftDeleteURLParams{
		Code:      input.Code,
		OwnerType: ownerType,
		OwnerID:   input.OwnerID,
	})
	if err != nil {
		s.Logger.Error("failed to delete url", "code", input.Code, "error", err)
		return ErrURLDelete
	}

	// Invalidate cache
	s.InvalidateURLCache(ctx, input.Code)

	s.Logger.Info("url deleted", "code", input.Code, "owner_id", input.OwnerID)
	return nil
}
