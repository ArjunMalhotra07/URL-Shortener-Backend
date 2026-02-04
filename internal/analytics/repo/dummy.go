package repo

import (
	"context"
	"time"

	db "url_shortner_backend/db/output"

	"github.com/jackc/pgx/v5/pgtype"
)

// AnalyticsRepoDummy returns simulated data for frontend testing
type AnalyticsRepoDummy struct{}

func NewAnalyticsRepoDummy() *AnalyticsRepoDummy {
	return &AnalyticsRepoDummy{}
}

func (r *AnalyticsRepoDummy) InsertClick(ctx context.Context, params db.InsertClickParams) error {
	return nil // No-op for dummy
}

func (r *AnalyticsRepoDummy) GetClicksByShortURLID(ctx context.Context, params db.GetClicksByShortURLIDParams) ([]db.Click, error) {
	now := time.Now()
	clicks := []db.Click{
		{
			ID:             1,
			ShortUrlID:     params.ShortUrlID,
			ClickedAt:      pgtype.Timestamptz{Time: now.Add(-1 * time.Hour), Valid: true},
			IpHash:         "abc123",
			Country:        pgtype.Text{String: "United States", Valid: true},
			City:           pgtype.Text{String: "New York", Valid: true},
			Browser:        pgtype.Text{String: "Chrome", Valid: true},
			Os:             pgtype.Text{String: "Windows 11", Valid: true},
			DeviceType:     pgtype.Text{String: "desktop", Valid: true},
			ReferrerDomain: pgtype.Text{String: "google.com", Valid: true},
			UtmSource:      pgtype.Text{String: "newsletter", Valid: true},
			UtmMedium:      pgtype.Text{String: "email", Valid: true},
			UtmCampaign:    pgtype.Text{String: "spring_sale", Valid: true},
			IsUnique:       pgtype.Bool{Bool: true, Valid: true},
			IsBot:          pgtype.Bool{Bool: false, Valid: true},
		},
		{
			ID:             2,
			ShortUrlID:     params.ShortUrlID,
			ClickedAt:      pgtype.Timestamptz{Time: now.Add(-2 * time.Hour), Valid: true},
			IpHash:         "def456",
			Country:        pgtype.Text{String: "United Kingdom", Valid: true},
			City:           pgtype.Text{String: "London", Valid: true},
			Browser:        pgtype.Text{String: "Safari", Valid: true},
			Os:             pgtype.Text{String: "macOS", Valid: true},
			DeviceType:     pgtype.Text{String: "desktop", Valid: true},
			ReferrerDomain: pgtype.Text{String: "twitter.com", Valid: true},
			IsUnique:       pgtype.Bool{Bool: true, Valid: true},
			IsBot:          pgtype.Bool{Bool: false, Valid: true},
		},
		{
			ID:             3,
			ShortUrlID:     params.ShortUrlID,
			ClickedAt:      pgtype.Timestamptz{Time: now.Add(-5 * time.Hour), Valid: true},
			IpHash:         "ghi789",
			Country:        pgtype.Text{String: "India", Valid: true},
			City:           pgtype.Text{String: "Mumbai", Valid: true},
			Browser:        pgtype.Text{String: "Chrome", Valid: true},
			Os:             pgtype.Text{String: "Android", Valid: true},
			DeviceType:     pgtype.Text{String: "mobile", Valid: true},
			ReferrerDomain: pgtype.Text{String: "facebook.com", Valid: true},
			IsUnique:       pgtype.Bool{Bool: true, Valid: true},
			IsBot:          pgtype.Bool{Bool: false, Valid: true},
		},
		{
			ID:             4,
			ShortUrlID:     params.ShortUrlID,
			ClickedAt:      pgtype.Timestamptz{Time: now.Add(-1 * 24 * time.Hour), Valid: true},
			IpHash:         "jkl012",
			Country:        pgtype.Text{String: "Germany", Valid: true},
			City:           pgtype.Text{String: "Berlin", Valid: true},
			Browser:        pgtype.Text{String: "Firefox", Valid: true},
			Os:             pgtype.Text{String: "Linux", Valid: true},
			DeviceType:     pgtype.Text{String: "desktop", Valid: true},
			ReferrerDomain: pgtype.Text{String: "linkedin.com", Valid: true},
			IsUnique:       pgtype.Bool{Bool: true, Valid: true},
			IsBot:          pgtype.Bool{Bool: false, Valid: true},
		},
		{
			ID:             5,
			ShortUrlID:     params.ShortUrlID,
			ClickedAt:      pgtype.Timestamptz{Time: now.Add(-2 * 24 * time.Hour), Valid: true},
			IpHash:         "mno345",
			Country:        pgtype.Text{String: "Japan", Valid: true},
			City:           pgtype.Text{String: "Tokyo", Valid: true},
			Browser:        pgtype.Text{String: "Edge", Valid: true},
			Os:             pgtype.Text{String: "iOS", Valid: true},
			DeviceType:     pgtype.Text{String: "tablet", Valid: true},
			ReferrerDomain: pgtype.Text{String: "reddit.com", Valid: true},
			IsUnique:       pgtype.Bool{Bool: true, Valid: true},
			IsBot:          pgtype.Bool{Bool: false, Valid: true},
		},
	}
	return clicks, nil
}

func (r *AnalyticsRepoDummy) CountClicksByShortURLID(ctx context.Context, shortUrlID int64) (int64, error) {
	return 1247, nil
}

func (r *AnalyticsRepoDummy) CountUniqueClicksByShortURLID(ctx context.Context, shortUrlID int64) (int64, error) {
	return 892, nil
}

func (r *AnalyticsRepoDummy) GetClicksByCountry(ctx context.Context, params db.GetClicksByCountryParams) ([]db.GetClicksByCountryRow, error) {
	return []db.GetClicksByCountryRow{
		{Country: pgtype.Text{String: "United States", Valid: true}, Clicks: 423},
		{Country: pgtype.Text{String: "India", Valid: true}, Clicks: 287},
		{Country: pgtype.Text{String: "United Kingdom", Valid: true}, Clicks: 156},
		{Country: pgtype.Text{String: "Germany", Valid: true}, Clicks: 98},
		{Country: pgtype.Text{String: "Canada", Valid: true}, Clicks: 87},
		{Country: pgtype.Text{String: "France", Valid: true}, Clicks: 64},
		{Country: pgtype.Text{String: "Japan", Valid: true}, Clicks: 52},
		{Country: pgtype.Text{String: "Australia", Valid: true}, Clicks: 43},
		{Country: pgtype.Text{String: "Brazil", Valid: true}, Clicks: 37},
	}, nil
}

func (r *AnalyticsRepoDummy) GetClicksByCity(ctx context.Context, params db.GetClicksByCityParams) ([]db.GetClicksByCityRow, error) {
	return []db.GetClicksByCityRow{
		{City: pgtype.Text{String: "New York", Valid: true}, Country: pgtype.Text{String: "United States", Valid: true}, Clicks: 142},
		{City: pgtype.Text{String: "Mumbai", Valid: true}, Country: pgtype.Text{String: "India", Valid: true}, Clicks: 98},
		{City: pgtype.Text{String: "London", Valid: true}, Country: pgtype.Text{String: "United Kingdom", Valid: true}, Clicks: 87},
		{City: pgtype.Text{String: "San Francisco", Valid: true}, Country: pgtype.Text{String: "United States", Valid: true}, Clicks: 76},
		{City: pgtype.Text{String: "Berlin", Valid: true}, Country: pgtype.Text{String: "Germany", Valid: true}, Clicks: 54},
		{City: pgtype.Text{String: "Bangalore", Valid: true}, Country: pgtype.Text{String: "India", Valid: true}, Clicks: 52},
		{City: pgtype.Text{String: "Toronto", Valid: true}, Country: pgtype.Text{String: "Canada", Valid: true}, Clicks: 48},
		{City: pgtype.Text{String: "Paris", Valid: true}, Country: pgtype.Text{String: "France", Valid: true}, Clicks: 41},
	}, nil
}

func (r *AnalyticsRepoDummy) GetClicksByDevice(ctx context.Context, params db.GetClicksByDeviceParams) ([]db.GetClicksByDeviceRow, error) {
	return []db.GetClicksByDeviceRow{
		{DeviceType: pgtype.Text{String: "desktop", Valid: true}, Clicks: 687},
		{DeviceType: pgtype.Text{String: "mobile", Valid: true}, Clicks: 498},
		{DeviceType: pgtype.Text{String: "tablet", Valid: true}, Clicks: 62},
	}, nil
}

func (r *AnalyticsRepoDummy) GetClicksByBrowser(ctx context.Context, params db.GetClicksByBrowserParams) ([]db.GetClicksByBrowserRow, error) {
	return []db.GetClicksByBrowserRow{
		{Browser: pgtype.Text{String: "Chrome", Valid: true}, Clicks: 612},
		{Browser: pgtype.Text{String: "Safari", Valid: true}, Clicks: 287},
		{Browser: pgtype.Text{String: "Firefox", Valid: true}, Clicks: 156},
		{Browser: pgtype.Text{String: "Edge", Valid: true}, Clicks: 98},
		{Browser: pgtype.Text{String: "Opera", Valid: true}, Clicks: 42},
		{Browser: pgtype.Text{String: "Other", Valid: true}, Clicks: 52},
	}, nil
}

func (r *AnalyticsRepoDummy) GetClicksByOS(ctx context.Context, params db.GetClicksByOSParams) ([]db.GetClicksByOSRow, error) {
	return []db.GetClicksByOSRow{
		{Os: pgtype.Text{String: "Windows", Valid: true}, Clicks: 423},
		{Os: pgtype.Text{String: "macOS", Valid: true}, Clicks: 287},
		{Os: pgtype.Text{String: "Android", Valid: true}, Clicks: 256},
		{Os: pgtype.Text{String: "iOS", Valid: true}, Clicks: 198},
		{Os: pgtype.Text{String: "Linux", Valid: true}, Clicks: 83},
	}, nil
}

func (r *AnalyticsRepoDummy) GetClicksTimeseries(ctx context.Context, params db.GetClicksTimeseriesParams) ([]db.GetClicksTimeseriesRow, error) {
	now := time.Now()
	return []db.GetClicksTimeseriesRow{
		{Date: pgtype.Timestamptz{Time: now.AddDate(0, 0, -6).Truncate(24 * time.Hour), Valid: true}, Clicks: 87},
		{Date: pgtype.Timestamptz{Time: now.AddDate(0, 0, -5).Truncate(24 * time.Hour), Valid: true}, Clicks: 124},
		{Date: pgtype.Timestamptz{Time: now.AddDate(0, 0, -4).Truncate(24 * time.Hour), Valid: true}, Clicks: 98},
		{Date: pgtype.Timestamptz{Time: now.AddDate(0, 0, -3).Truncate(24 * time.Hour), Valid: true}, Clicks: 156},
		{Date: pgtype.Timestamptz{Time: now.AddDate(0, 0, -2).Truncate(24 * time.Hour), Valid: true}, Clicks: 203},
		{Date: pgtype.Timestamptz{Time: now.AddDate(0, 0, -1).Truncate(24 * time.Hour), Valid: true}, Clicks: 178},
		{Date: pgtype.Timestamptz{Time: now.Truncate(24 * time.Hour), Valid: true}, Clicks: 142},
	}, nil
}

func (r *AnalyticsRepoDummy) GetClicksTimeseriesHourly(ctx context.Context, params db.GetClicksTimeseriesHourlyParams) ([]db.GetClicksTimeseriesHourlyRow, error) {
	now := time.Now()
	rows := []db.GetClicksTimeseriesHourlyRow{}
	for i := 23; i >= 0; i-- {
		clicks := int64(10 + (i*7)%50) // Varying pattern
		rows = append(rows, db.GetClicksTimeseriesHourlyRow{
			Date:   pgtype.Timestamptz{Time: now.Add(-time.Duration(i) * time.Hour).Truncate(time.Hour), Valid: true},
			Clicks: clicks,
		})
	}
	return rows, nil
}

func (r *AnalyticsRepoDummy) GetTopReferrers(ctx context.Context, params db.GetTopReferrersParams) ([]db.GetTopReferrersRow, error) {
	return []db.GetTopReferrersRow{
		{ReferrerDomain: pgtype.Text{String: "google.com", Valid: true}, Clicks: 312},
		{ReferrerDomain: pgtype.Text{String: "twitter.com", Valid: true}, Clicks: 187},
		{ReferrerDomain: pgtype.Text{String: "facebook.com", Valid: true}, Clicks: 143},
		{ReferrerDomain: pgtype.Text{String: "linkedin.com", Valid: true}, Clicks: 98},
		{ReferrerDomain: pgtype.Text{String: "reddit.com", Valid: true}, Clicks: 76},
		{ReferrerDomain: pgtype.Text{String: "github.com", Valid: true}, Clicks: 54},
		{ReferrerDomain: pgtype.Text{String: "youtube.com", Valid: true}, Clicks: 42},
	}, nil
}

func (r *AnalyticsRepoDummy) GetTopCampaigns(ctx context.Context, params db.GetTopCampaignsParams) ([]db.GetTopCampaignsRow, error) {
	return []db.GetTopCampaignsRow{
		{UtmCampaign: pgtype.Text{String: "spring_sale", Valid: true}, Clicks: 234},
		{UtmCampaign: pgtype.Text{String: "product_launch", Valid: true}, Clicks: 156},
		{UtmCampaign: pgtype.Text{String: "newsletter_feb", Valid: true}, Clicks: 98},
		{UtmCampaign: pgtype.Text{String: "social_promo", Valid: true}, Clicks: 67},
	}, nil
}

func (r *AnalyticsRepoDummy) GetClicksSummary(ctx context.Context, params db.GetClicksSummaryParams) (db.GetClicksSummaryRow, error) {
	return db.GetClicksSummaryRow{
		TotalClicks:  1247,
		UniqueClicks: 892,
		BotClicks:    23,
	}, nil
}

func (r *AnalyticsRepoDummy) GetClicksSummaryAllTime(ctx context.Context, shortUrlID int64) (db.GetClicksSummaryAllTimeRow, error) {
	return db.GetClicksSummaryAllTimeRow{
		TotalClicks:  4521,
		UniqueClicks: 3187,
		BotClicks:    89,
	}, nil
}
