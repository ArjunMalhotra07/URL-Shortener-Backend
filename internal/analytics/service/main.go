package service

import (
	"context"
	"time"

	db "url_shortner_backend/db/output"
	"url_shortner_backend/internal/analytics/repo"
	shorturlrepo "url_shortner_backend/internal/shorturl/repo"
	"url_shortner_backend/pkg/geoip"
	"url_shortner_backend/pkg/redis"

	"github.com/rs/zerolog"
)

type AnalyticsSvc interface {
	RecordClick(ctx context.Context, input RecordClickInput) error
	GetSummary(ctx context.Context, input GetSummaryInput) (GetSummaryOutput, error)
	GetClicks(ctx context.Context, input GetClicksInput) (GetClicksOutput, error)
	GetTimeseries(ctx context.Context, input GetTimeseriesInput) (GetTimeseriesOutput, error)
	GetGeoBreakdown(ctx context.Context, input GetGeoInput) (GetGeoOutput, error)
	GetDeviceBreakdown(ctx context.Context, input GetDeviceInput) (GetDeviceOutput, error)
}

type AnalyticsSvcImp struct {
	Logger       zerolog.Logger
	Repo         repo.AnalyticsRepository
	ShortURLRepo shorturlrepo.ShortURLRepository
	GeoIP        geoip.GeoIPLookup
	Redis        *redis.RedisClient
}

type AnalyticsSvcParams struct {
	Logger       zerolog.Logger
	Repo         repo.AnalyticsRepository
	ShortURLRepo shorturlrepo.ShortURLRepository
	GeoIP        geoip.GeoIPLookup
	Redis        *redis.RedisClient
}

func NewAnalyticsSvcImp(p AnalyticsSvcParams) *AnalyticsSvcImp {
	return &AnalyticsSvcImp{
		Logger:       p.Logger,
		Repo:         p.Repo,
		ShortURLRepo: p.ShortURLRepo,
		GeoIP:        p.GeoIP,
		Redis:        p.Redis,
	}
}

// Common types

type CountryStats struct {
	Country string `json:"country"`
	Clicks  int64  `json:"clicks"`
}

type DeviceTypeStats struct {
	DeviceType string `json:"device_type"`
	Clicks     int64  `json:"clicks"`
}

// Helper to check if time range is "all time" (zero start time)
func isAllTime(start time.Time) bool {
	return start.IsZero()
}

// Helper to get start time from time range - returns zero for "all"
func getStartTime(start time.Time) time.Time {
	if start.IsZero() {
		return time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	return start
}

// Helper to get short URL and validate ownership
func (s *AnalyticsSvcImp) getAndValidateURL(ctx context.Context, code, ownerType, ownerID string) (db.ShortUrl, error) {
	shortURL, err := s.ShortURLRepo.GetShortURLByCode(ctx, code)
	if err != nil {
		return db.ShortUrl{}, ErrURLNotFound
	}

	if string(shortURL.OwnerType) != ownerType || shortURL.OwnerID != ownerID {
		return db.ShortUrl{}, ErrURLNotOwned
	}

	return shortURL, nil
}
