package service

import (
	"context"

	"github.com/jackc/pgx/v5"

	db "url_shortner_backend/db/output"
)

type UpdateLongURLInput struct {
	Code      string
	LongURL   string
	OwnerType string
	OwnerID   string
}

func (s *ShortURLSvcImp) UpdateLongURL(ctx context.Context, input UpdateLongURLInput) error {
	if input.Code == "" {
		return ErrInvalidCode
	}
	if input.OwnerID == "" {
		return ErrInvalidOwner
	}

	longURL := normalizeURL(input.LongURL)
	if err := validateURL(longURL); err != nil {
		s.Logger.Error("invalid url provided", "url", longURL, "error", err)
		return ErrInvalidURL
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
			s.Logger.Info("update attempt on non-owned url", "code", input.Code, "owner_id", input.OwnerID)
			return ErrURLNotOwned
		}
		s.Logger.Error("failed to verify url ownership", "code", input.Code, "error", err)
		return ErrURLUpdate
	}

	err = s.Repo.UpdateLongURL(ctx, db.UpdateLongURLParams{
		Code:      input.Code,
		OwnerType: ownerType,
		OwnerID:   input.OwnerID,
		LongUrl:   longURL,
	})
	if err != nil {
		s.Logger.Error("failed to update long url", "code", input.Code, "error", err)
		return ErrURLUpdate
	}

	s.Logger.Info("long url updated", "code", input.Code, "owner_id", input.OwnerID)
	return nil
}
