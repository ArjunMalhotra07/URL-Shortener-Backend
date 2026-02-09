package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type GetLongURLInput struct {
	Code string
}
type GetLongURLOutput struct {
	ID      int64
	LongURL string
}

// cachedURL is stored in Redis
type cachedURL struct {
	ID      int64  `json:"id"`
	LongURL string `json:"long_url"`
}

const urlCacheTTL = 10 * time.Minute

func (s *ShortURLSvcImp) GetLongURL(ctx context.Context, input GetLongURLInput) (GetLongURLOutput, error) {
	if input.Code == "" {
		s.Logger.Error().Msg("empty code provided")
		return GetLongURLOutput{}, ErrInvalidCode
	}

	cacheKey := fmt.Sprintf("url:%s", input.Code)

	// Try Redis cache first
	if s.Redis != nil {
		cached, err := s.Redis.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			var cu cachedURL
			if json.Unmarshal([]byte(cached), &cu) == nil {
				// Extend TTL on hit (sliding expiration)
				_ = s.Redis.Expire(ctx, cacheKey, urlCacheTTL)
				s.Logger.Debug().Str("code", input.Code).Msg("cache hit")
				return GetLongURLOutput{
					ID:      cu.ID,
					LongURL: cu.LongURL,
				}, nil
			}
		}
	}

	// Cache miss - query DB
	shortURL, err := s.Repo.GetShortURLByCode(ctx, input.Code)
	if err != nil {
		s.Logger.Err(err).Str("code", input.Code).Msg("failed to get short url")
		return GetLongURLOutput{}, ErrURLNotFound
	}

	// Check if URL is active
	if !shortURL.IsActive {
		s.Logger.Info().Str("code", input.Code).Msg("short url is inactive")
		return GetLongURLOutput{}, ErrURLInactive
	}

	// Check if URL has expired
	if shortURL.ExpiresAt.Valid && shortURL.ExpiresAt.Time.Before(time.Now()) {
		s.Logger.Info().Str("code", input.Code).Msg("short url has expired")
		return GetLongURLOutput{}, ErrURLExpired
	}

	// Cache the result (only cache active, non-expired URLs)
	if s.Redis != nil {
		cu := cachedURL{ID: shortURL.ID, LongURL: shortURL.LongUrl}
		if data, err := json.Marshal(cu); err == nil {
			ttl := urlCacheTTL
			// If URL has expiry, use shorter TTL
			if shortURL.ExpiresAt.Valid {
				timeUntilExpiry := time.Until(shortURL.ExpiresAt.Time)
				if timeUntilExpiry < ttl {
					ttl = timeUntilExpiry
				}
			}
			if ttl > 0 {
				_ = s.Redis.Set(ctx, cacheKey, string(data), ttl)
				s.Logger.Debug().Str("code", input.Code).Dur("ttl", ttl).Msg("cached url")
			}
		}
	}
	return GetLongURLOutput{ID: shortURL.ID, LongURL: shortURL.LongUrl}, nil
}

// InvalidateURLCache removes a URL from cache (call when URL is updated/toggled/deleted)
func (s *ShortURLSvcImp) InvalidateURLCache(ctx context.Context, code string) {
	if s.Redis != nil {
		cacheKey := fmt.Sprintf("url:%s", code)
		_ = s.Redis.Del(ctx, cacheKey)
		s.Logger.Debug().Str("code", code).Msg("invalidated cache")
	}
}
