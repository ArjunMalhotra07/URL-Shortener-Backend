package service

import (
	"context"
	"time"
	db "url_shortner_backend/db/output"
)

type GetMyURLsInput struct {
	OwnerType string
	OwnerID   string
	Limit     int32
	Offset    int32
}
type GetMyURLsOutput struct {
	URLs  []ShortURLItem
	Total int64
}
type ShortURLItem struct {
	ID        int64
	Code      string
	LongURL   string
	IsActive  bool
	ExpiresAt *time.Time
	IsExpired bool
	CreatedAt time.Time
}

func (s *ShortURLSvcImp) GetMyURLs(ctx context.Context, input GetMyURLsInput) (GetMyURLsOutput, error) {
	if input.OwnerID == "" {
		s.Logger.Error("empty owner id provided")
		return GetMyURLsOutput{}, ErrInvalidOwner
	}

	ownerType := db.OwnerTypeEnumAnonymous
	if input.OwnerType == "user" {
		ownerType = db.OwnerTypeEnumUser
	}

	urls, err := s.Repo.GetShortURLsByOwner(ctx, db.GetShortURLsByOwnerParams{
		OwnerType: ownerType,
		OwnerID:   input.OwnerID,
		Limit:     input.Limit,
		Offset:    input.Offset,
	})
	if err != nil {
		s.Logger.Error("failed to get urls by owner", "owner_id", input.OwnerID, "error", err)
		return GetMyURLsOutput{}, ErrURLFetch
	}

	total, err := s.Repo.CountURLsByOwner(ctx, db.CountURLsByOwnerParams{
		OwnerType: ownerType,
		OwnerID:   input.OwnerID,
	})
	if err != nil {
		s.Logger.Error("failed to count urls by owner", "owner_id", input.OwnerID, "error", err)
		return GetMyURLsOutput{}, ErrURLFetch
	}

	items := make([]ShortURLItem, len(urls))
	now := time.Now()
	for i, u := range urls {
		var expiresAt *time.Time
		var isExpired bool
		if u.ExpiresAt.Valid {
			expiresAt = &u.ExpiresAt.Time
			isExpired = u.ExpiresAt.Time.Before(now)
		}
		items[i] = ShortURLItem{
			ID:        u.ID,
			Code:      u.Code,
			LongURL:   u.LongUrl,
			IsActive:  u.IsActive,
			ExpiresAt: expiresAt,
			IsExpired: isExpired,
			CreatedAt: u.CreatedAt.Time,
		}
	}

	s.Logger.Info("urls retrieved", "owner_id", input.OwnerID, "count", len(items), "total", total)

	return GetMyURLsOutput{URLs: items, Total: total}, nil
}
