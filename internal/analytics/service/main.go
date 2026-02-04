package service

import (
	"context"
	"time"

	db "url_shortner_backend/db/output"
	"url_shortner_backend/internal/analytics/repo"
	shorturlrepo "url_shortner_backend/internal/shorturl/repo"
	"url_shortner_backend/pkg/geoip"
	"url_shortner_backend/pkg/logger"
	"url_shortner_backend/pkg/redis"
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
	Logger       logger.Logger
	Repo         repo.AnalyticsRepository
	ShortURLRepo shorturlrepo.ShortURLRepository
	GeoIP        geoip.GeoIPLookup
	Redis        *redis.RedisClient
}

type AnalyticsSvcParams struct {
	Logger       logger.Logger
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

// Helper function to parse time range
func parseTimeRange(timeRange string) time.Time {
	now := time.Now()
	switch timeRange {
	case "24h":
		return now.Add(-24 * time.Hour)
	case "7d":
		return now.Add(-7 * 24 * time.Hour)
	case "30d":
		return now.Add(-30 * 24 * time.Hour)
	case "90d":
		return now.Add(-90 * 24 * time.Hour)
	default:
		return time.Time{} // zero time for "all"
	}
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
