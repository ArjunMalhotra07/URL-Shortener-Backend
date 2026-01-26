package service

import (
	"context"
	"time"
)

type GetLongURLInput struct {
	Code string
}
type GetLongURLOutput struct {
	LongURL string
}

func (s *ShortURLSvcImp) GetLongURL(ctx context.Context, input GetLongURLInput) (GetLongURLOutput, error) {
	if input.Code == "" {
		s.Logger.Error("empty code provided")
		return GetLongURLOutput{}, ErrInvalidCode
	}

	shortURL, err := s.Repo.GetShortURLByCode(ctx, input.Code)
	if err != nil {
		s.Logger.Error("failed to get short url", "code", input.Code, "error", err)
		return GetLongURLOutput{}, ErrURLNotFound
	}

	// Check if URL is active
	if !shortURL.IsActive {
		s.Logger.Info("short url is inactive", "code", input.Code)
		return GetLongURLOutput{}, ErrURLInactive
	}

	// Check if URL has expired
	if shortURL.ExpiresAt.Valid && shortURL.ExpiresAt.Time.Before(time.Now()) {
		s.Logger.Info("short url has expired", "code", input.Code)
		return GetLongURLOutput{}, ErrURLExpired
	}

	s.Logger.Info("long url retrieved", "code", input.Code)

	return GetLongURLOutput{
		LongURL: shortURL.LongUrl,
	}, nil
}
