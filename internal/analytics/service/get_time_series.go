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
	Start     time.Time // zero time means all time
	End       time.Time
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

	clickedAtParam := pgtype.Timestamptz{Time: getStartTime(input.Start), Valid: true}

	var data []TimeseriesPoint

	if input.Interval == "hour" {
		rows, err := s.Repo.GetClicksTimeseriesHourly(ctx, db.GetClicksTimeseriesHourlyParams{
			ShortUrlID: shortURL.ID,
			ClickedAt:  clickedAtParam,
		})
		if err != nil {
			s.Logger.Err(err).Msg("failed to get hourly timeseries")
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
		rows, err := s.Repo.GetClicksTimeseries(ctx, db.GetClicksTimeseriesParams{ShortUrlID: shortURL.ID, ClickedAt: clickedAtParam})
		if err != nil {
			s.Logger.Err(err).Msg("failed to get daily timeseries")
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
