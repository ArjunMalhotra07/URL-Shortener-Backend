package service

import (
	"context"
	"net/url"
	"strings"

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

	// Validate and normalize the new long URL
	longURL := strings.TrimSpace(input.LongURL)
	if longURL == "" {
		return ErrInvalidURL
	}

	// Add https:// if no scheme provided
	if !strings.HasPrefix(longURL, "http://") && !strings.HasPrefix(longURL, "https://") {
		longURL = "https://" + longURL
	}

	parsedURL, err := url.ParseRequestURI(longURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return ErrInvalidURL
	}

	ownerType := db.OwnerTypeEnumAnonymous
	if input.OwnerType == "user" {
		ownerType = db.OwnerTypeEnumUser
	}

	// Verify ownership
	_, err = s.Repo.GetURLByCodeAndOwner(ctx, db.GetURLByCodeAndOwnerParams{
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
