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

// Input/Output types

type RecordClickInput struct {
	ShortURLID  int64
	IPAddress   string
	UserAgent   string
	Referrer    string
	UTMSource   string
	UTMMedium   string
	UTMCampaign string
}

type GetSummaryInput struct {
	Code      string
	OwnerType string
	OwnerID   string
	TimeRange string // "24h", "7d", "30d", "all"
}

type GetSummaryOutput struct {
	TotalClicks  int64             `json:"total_clicks"`
	UniqueClicks int64             `json:"unique_clicks"`
	BotClicks    int64             `json:"bot_clicks"`
	TopCountries []CountryStats    `json:"top_countries"`
	TopReferrers []ReferrerStats   `json:"top_referrers"`
	DeviceStats  []DeviceTypeStats `json:"device_stats"`
}

type CountryStats struct {
	Country string `json:"country"`
	Clicks  int64  `json:"clicks"`
}

type ReferrerStats struct {
	Domain string `json:"domain"`
	Clicks int64  `json:"clicks"`
}

type DeviceTypeStats struct {
	DeviceType string `json:"device_type"`
	Clicks     int64  `json:"clicks"`
}

type GetClicksInput struct {
	Code      string
	OwnerType string
	OwnerID   string
	Limit     int32
	Offset    int32
}

type GetClicksOutput struct {
	Clicks     []ClickRecord `json:"clicks"`
	TotalCount int64         `json:"total_count"`
}

type ClickRecord struct {
	ID             int64     `json:"id"`
	ClickedAt      time.Time `json:"clicked_at"`
	Country        string    `json:"country,omitempty"`
	City           string    `json:"city,omitempty"`
	Browser        string    `json:"browser,omitempty"`
	OS             string    `json:"os,omitempty"`
	DeviceType     string    `json:"device_type,omitempty"`
	Referrer       string    `json:"referrer,omitempty"`
	ReferrerDomain string    `json:"referrer_domain,omitempty"`
	UTMSource      string    `json:"utm_source,omitempty"`
	UTMMedium      string    `json:"utm_medium,omitempty"`
	UTMCampaign    string    `json:"utm_campaign,omitempty"`
	IsUnique       bool      `json:"is_unique"`
	IsBot          bool      `json:"is_bot"`
}

type GetTimeseriesInput struct {
	Code      string
	OwnerType string
	OwnerID   string
	TimeRange string // "24h", "7d", "30d"
	Interval  string // "hour", "day"
}

type GetTimeseriesOutput struct {
	Data []TimeseriesPoint `json:"data"`
}

type TimeseriesPoint struct {
	Date   time.Time `json:"date"`
	Clicks int64     `json:"clicks"`
}

type GetGeoInput struct {
	Code      string
	OwnerType string
	OwnerID   string
	TimeRange string
	Limit     int32
}

type GetGeoOutput struct {
	Countries []CountryStats `json:"countries"`
	Cities    []CityStats    `json:"cities"`
}

type CityStats struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Clicks  int64  `json:"clicks"`
}

type GetDeviceInput struct {
	Code      string
	OwnerType string
	OwnerID   string
	TimeRange string
}

type GetDeviceOutput struct {
	DeviceTypes []DeviceTypeStats `json:"device_types"`
	Browsers    []BrowserStats    `json:"browsers"`
	OS          []OSStats         `json:"os"`
}

type BrowserStats struct {
	Browser string `json:"browser"`
	Clicks  int64  `json:"clicks"`
}

type OSStats struct {
	OS     string `json:"os"`
	Clicks int64  `json:"clicks"`
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
