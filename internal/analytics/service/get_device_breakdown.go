package service

import (
	"context"
	"time"
	db "url_shortner_backend/db/output"

	"github.com/jackc/pgx/v5/pgtype"
)

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
type OSStats struct {
	OS     string `json:"os"`
	Clicks int64  `json:"clicks"`
}

type BrowserStats struct {
	Browser string `json:"browser"`
	Clicks  int64  `json:"clicks"`
}

func (s *AnalyticsSvcImp) GetDeviceBreakdown(ctx context.Context, input GetDeviceInput) (GetDeviceOutput, error) {
	// Validate ownership
	shortURL, err := s.getAndValidateURL(ctx, input.Code, input.OwnerType, input.OwnerID)
	if err != nil {
		return GetDeviceOutput{}, err
	}

	since := parseTimeRange(input.TimeRange)
	clickedAtParam := pgtype.Timestamptz{Valid: true}
	if !since.IsZero() {
		clickedAtParam.Time = since
	} else {
		clickedAtParam.Time = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	// Get device types
	devices, err := s.Repo.GetClicksByDevice(ctx, db.GetClicksByDeviceParams{
		ShortUrlID: shortURL.ID,
		ClickedAt:  clickedAtParam,
	})
	if err != nil {
		s.Logger.Error("failed to get devices", "error", err)
		return GetDeviceOutput{}, ErrAnalyticsFetch
	}

	deviceStats := make([]DeviceTypeStats, 0, len(devices))
	for _, d := range devices {
		if d.DeviceType.Valid {
			deviceStats = append(deviceStats, DeviceTypeStats{
				DeviceType: d.DeviceType.String,
				Clicks:     d.Clicks,
			})
		}
	}

	// Get browsers
	browsers, err := s.Repo.GetClicksByBrowser(ctx, db.GetClicksByBrowserParams{
		ShortUrlID: shortURL.ID,
		ClickedAt:  clickedAtParam,
	})
	if err != nil {
		s.Logger.Error("failed to get browsers", "error", err)
		return GetDeviceOutput{}, ErrAnalyticsFetch
	}

	browserStats := make([]BrowserStats, 0, len(browsers))
	for _, b := range browsers {
		if b.Browser.Valid {
			browserStats = append(browserStats, BrowserStats{
				Browser: b.Browser.String,
				Clicks:  b.Clicks,
			})
		}
	}

	// Get OS
	osData, err := s.Repo.GetClicksByOS(ctx, db.GetClicksByOSParams{
		ShortUrlID: shortURL.ID,
		ClickedAt:  clickedAtParam,
	})
	if err != nil {
		s.Logger.Error("failed to get OS", "error", err)
		return GetDeviceOutput{}, ErrAnalyticsFetch
	}

	osStats := make([]OSStats, 0, len(osData))
	for _, o := range osData {
		if o.Os.Valid {
			osStats = append(osStats, OSStats{
				OS:     o.Os.String,
				Clicks: o.Clicks,
			})
		}
	}

	return GetDeviceOutput{
		DeviceTypes: deviceStats,
		Browsers:    browserStats,
		OS:          osStats,
	}, nil
}
