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
	Start     time.Time // zero time means all time
	End       time.Time
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

	clickedAtParam := pgtype.Timestamptz{Time: getStartTime(input.Start), Valid: true}

	// Get device types
	devices, err := s.Repo.GetClicksByDevice(ctx, db.GetClicksByDeviceParams{
		ShortUrlID: shortURL.ID,
		ClickedAt:  clickedAtParam,
	})
	if err != nil {
		s.Logger.Err(err).Msg("failed to get devices")
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
		s.Logger.Err(err).Msg("failed to get browsers")
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
		s.Logger.Err(err).Msg("failed to get OS")
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
