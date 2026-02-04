package service

import (
	"context"
	"time"
	db "url_shortner_backend/db/output"

	"github.com/jackc/pgx/v5/pgtype"
)

type TimeseriesPoint struct {
	Date   time.Time `json:"date"`
	Clicks int64     `json:"clicks"`
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
