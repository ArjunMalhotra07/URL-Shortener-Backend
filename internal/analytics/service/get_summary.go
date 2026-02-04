package service

import (
	"context"
	"time"

	db "url_shortner_backend/db/output"

	"github.com/jackc/pgx/v5/pgtype"
)

func (s *AnalyticsSvcImp) GetSummary(ctx context.Context, input GetSummaryInput) (GetSummaryOutput, error) {
	// Validate ownership
	shortURL, err := s.getAndValidateURL(ctx, input.Code, input.OwnerType, input.OwnerID)
	if err != nil {
		return GetSummaryOutput{}, err
	}

	since := parseTimeRange(input.TimeRange)

	var summary db.GetClicksSummaryRow
	if since.IsZero() {
		// All time
		allTimeSummary, err := s.Repo.GetClicksSummaryAllTime(ctx, shortURL.ID)
		if err != nil {
			s.Logger.Error("failed to get all-time summary", "error", err)
			return GetSummaryOutput{}, ErrAnalyticsFetch
		}
		summary = db.GetClicksSummaryRow{
			TotalClicks:  allTimeSummary.TotalClicks,
			UniqueClicks: allTimeSummary.UniqueClicks,
			BotClicks:    allTimeSummary.BotClicks,
		}
	} else {
		summary, err = s.Repo.GetClicksSummary(ctx, db.GetClicksSummaryParams{
			ShortUrlID: shortURL.ID,
			ClickedAt:  pgtype.Timestamptz{Time: since, Valid: true},
		})
		if err != nil {
			s.Logger.Error("failed to get summary", "error", err)
			return GetSummaryOutput{}, ErrAnalyticsFetch
		}
	}

	// Get top countries
	clickedAtParam := pgtype.Timestamptz{Valid: true}
	if !since.IsZero() {
		clickedAtParam.Time = since
	} else {
		clickedAtParam.Time = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	countries, err := s.Repo.GetClicksByCountry(ctx, db.GetClicksByCountryParams{
		ShortUrlID: shortURL.ID,
		ClickedAt:  clickedAtParam,
		Limit:      5,
	})
	if err != nil {
		s.Logger.Error("failed to get countries", "error", err)
	}

	topCountries := make([]CountryStats, 0, len(countries))
	for _, c := range countries {
		if c.Country.Valid {
			topCountries = append(topCountries, CountryStats{
				Country: c.Country.String,
				Clicks:  c.Clicks,
			})
		}
	}

	// Get top referrers
	referrers, err := s.Repo.GetTopReferrers(ctx, db.GetTopReferrersParams{
		ShortUrlID: shortURL.ID,
		ClickedAt:  clickedAtParam,
		Limit:      5,
	})
	if err != nil {
		s.Logger.Error("failed to get referrers", "error", err)
	}

	topReferrers := make([]ReferrerStats, 0, len(referrers))
	for _, r := range referrers {
		if r.ReferrerDomain.Valid {
			topReferrers = append(topReferrers, ReferrerStats{
				Domain: r.ReferrerDomain.String,
				Clicks: r.Clicks,
			})
		}
	}

	// Get device stats
	devices, err := s.Repo.GetClicksByDevice(ctx, db.GetClicksByDeviceParams{
		ShortUrlID: shortURL.ID,
		ClickedAt:  clickedAtParam,
	})
	if err != nil {
		s.Logger.Error("failed to get devices", "error", err)
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

	return GetSummaryOutput{
		TotalClicks:  summary.TotalClicks,
		UniqueClicks: summary.UniqueClicks,
		BotClicks:    summary.BotClicks,
		TopCountries: topCountries,
		TopReferrers: topReferrers,
		DeviceStats:  deviceStats,
	}, nil
}

func (s *AnalyticsSvcImp) GetClicks(ctx context.Context, input GetClicksInput) (GetClicksOutput, error) {
	// Validate ownership
	shortURL, err := s.getAndValidateURL(ctx, input.Code, input.OwnerType, input.OwnerID)
	if err != nil {
		return GetClicksOutput{}, err
	}

	// Set defaults
	limit := input.Limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	clicks, err := s.Repo.GetClicksByShortURLID(ctx, db.GetClicksByShortURLIDParams{
		ShortUrlID: shortURL.ID,
		Limit:      limit,
		Offset:     input.Offset,
	})
	if err != nil {
		s.Logger.Error("failed to get clicks", "error", err)
		return GetClicksOutput{}, ErrAnalyticsFetch
	}

	totalCount, err := s.Repo.CountClicksByShortURLID(ctx, shortURL.ID)
	if err != nil {
		s.Logger.Error("failed to count clicks", "error", err)
		return GetClicksOutput{}, ErrAnalyticsFetch
	}

	records := make([]ClickRecord, 0, len(clicks))
	for _, c := range clicks {
		records = append(records, ClickRecord{
			ID:             c.ID,
			ClickedAt:      c.ClickedAt.Time,
			Country:        c.Country.String,
			City:           c.City.String,
			Browser:        c.Browser.String,
			OS:             c.Os.String,
			DeviceType:     c.DeviceType.String,
			Referrer:       c.Referrer.String,
			ReferrerDomain: c.ReferrerDomain.String,
			UTMSource:      c.UtmSource.String,
			UTMMedium:      c.UtmMedium.String,
			UTMCampaign:    c.UtmCampaign.String,
			IsUnique:       c.IsUnique.Bool,
			IsBot:          c.IsBot.Bool,
		})
	}

	return GetClicksOutput{
		Clicks:     records,
		TotalCount: totalCount,
	}, nil
}

func (s *AnalyticsSvcImp) GetTimeseries(ctx context.Context, input GetTimeseriesInput) (GetTimeseriesOutput, error) {
	// Validate ownership
	shortURL, err := s.getAndValidateURL(ctx, input.Code, input.OwnerType, input.OwnerID)
	if err != nil {
		return GetTimeseriesOutput{}, err
	}

	since := parseTimeRange(input.TimeRange)
	if since.IsZero() {
		since = time.Now().Add(-30 * 24 * time.Hour) // default to 30 days
	}

	clickedAtParam := pgtype.Timestamptz{Time: since, Valid: true}

	var data []TimeseriesPoint

	if input.Interval == "hour" {
		rows, err := s.Repo.GetClicksTimeseriesHourly(ctx, db.GetClicksTimeseriesHourlyParams{
			ShortUrlID: shortURL.ID,
			ClickedAt:  clickedAtParam,
		})
		if err != nil {
			s.Logger.Error("failed to get hourly timeseries", "error", err)
			return GetTimeseriesOutput{}, ErrAnalyticsFetch
		}

		data = make([]TimeseriesPoint, 0, len(rows))
		for _, r := range rows {
			data = append(data, TimeseriesPoint{
				Date:   r.Date.Time,
				Clicks: r.Clicks,
			})
		}
	} else {
		rows, err := s.Repo.GetClicksTimeseries(ctx, db.GetClicksTimeseriesParams{
			ShortUrlID: shortURL.ID,
			ClickedAt:  clickedAtParam,
		})
		if err != nil {
			s.Logger.Error("failed to get daily timeseries", "error", err)
			return GetTimeseriesOutput{}, ErrAnalyticsFetch
		}

		data = make([]TimeseriesPoint, 0, len(rows))
		for _, r := range rows {
			data = append(data, TimeseriesPoint{
				Date:   r.Date.Time,
				Clicks: r.Clicks,
			})
		}
	}

	return GetTimeseriesOutput{Data: data}, nil
}

func (s *AnalyticsSvcImp) GetGeoBreakdown(ctx context.Context, input GetGeoInput) (GetGeoOutput, error) {
	// Validate ownership
	shortURL, err := s.getAndValidateURL(ctx, input.Code, input.OwnerType, input.OwnerID)
	if err != nil {
		return GetGeoOutput{}, err
	}

	since := parseTimeRange(input.TimeRange)
	clickedAtParam := pgtype.Timestamptz{Valid: true}
	if !since.IsZero() {
		clickedAtParam.Time = since
	} else {
		clickedAtParam.Time = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	limit := input.Limit
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	// Get countries
	countries, err := s.Repo.GetClicksByCountry(ctx, db.GetClicksByCountryParams{
		ShortUrlID: shortURL.ID,
		ClickedAt:  clickedAtParam,
		Limit:      limit,
	})
	if err != nil {
		s.Logger.Error("failed to get countries", "error", err)
		return GetGeoOutput{}, ErrAnalyticsFetch
	}

	countryStats := make([]CountryStats, 0, len(countries))
	for _, c := range countries {
		if c.Country.Valid {
			countryStats = append(countryStats, CountryStats{
				Country: c.Country.String,
				Clicks:  c.Clicks,
			})
		}
	}

	// Get cities
	cities, err := s.Repo.GetClicksByCity(ctx, db.GetClicksByCityParams{
		ShortUrlID: shortURL.ID,
		ClickedAt:  clickedAtParam,
		Limit:      limit,
	})
	if err != nil {
		s.Logger.Error("failed to get cities", "error", err)
		return GetGeoOutput{}, ErrAnalyticsFetch
	}

	cityStats := make([]CityStats, 0, len(cities))
	for _, c := range cities {
		if c.City.Valid {
			cityStats = append(cityStats, CityStats{
				City:    c.City.String,
				Country: c.Country.String,
				Clicks:  c.Clicks,
			})
		}
	}

	return GetGeoOutput{
		Countries: countryStats,
		Cities:    cityStats,
	}, nil
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
