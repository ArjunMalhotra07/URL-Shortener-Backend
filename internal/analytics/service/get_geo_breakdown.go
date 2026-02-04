package service

import (
	"context"
	"time"
	db "url_shortner_backend/db/output"

	"github.com/jackc/pgx/v5/pgtype"
)

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
