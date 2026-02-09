package service

import (
	"context"
	"time"

	db "url_shortner_backend/db/output"

	"github.com/jackc/pgx/v5/pgtype"
)

type GetSummaryInput struct {
	Code      string
	OwnerType string
	OwnerID   string
	Start     time.Time // zero time means all time
	End       time.Time
}

type GetSummaryOutput struct {
	TotalClicks  int64           `json:"total_clicks"`
	UniqueClicks int64           `json:"unique_clicks"`
	BotClicks    int64           `json:"bot_clicks"`
	TopReferrers []ReferrerStats `json:"top_referrers"`
}

type ReferrerStats struct {
	Domain string `json:"domain"`
	Clicks int64  `json:"clicks"`
}

func (s *AnalyticsSvcImp) GetSummary(ctx context.Context, input GetSummaryInput) (GetSummaryOutput, error) {
	// Validate ownership
	shortURL, err := s.getAndValidateURL(ctx, input.Code, input.OwnerType, input.OwnerID)
	if err != nil {
		return GetSummaryOutput{}, err
	}

	var summary db.GetClicksSummaryRow
	if isAllTime(input.Start) {
		// All time
		allTimeSummary, err := s.Repo.GetClicksSummaryAllTime(ctx, shortURL.ID)
		if err != nil {
			s.Logger.Err(err).Msg("failed to get all-time summary")
			return GetSummaryOutput{}, ErrAnalyticsFetch
		}
		summary = db.GetClicksSummaryRow{TotalClicks: allTimeSummary.TotalClicks, UniqueClicks: allTimeSummary.UniqueClicks, BotClicks: allTimeSummary.BotClicks}
	} else {
		summary, err = s.Repo.GetClicksSummary(ctx, db.GetClicksSummaryParams{ShortUrlID: shortURL.ID, ClickedAt: pgtype.Timestamptz{Time: input.Start, Valid: true}})
		if err != nil {
			s.Logger.Err(err).Msg("failed to get summary")
			return GetSummaryOutput{}, ErrAnalyticsFetch
		}
	}

	// Get top referrers
	clickedAtParam := pgtype.Timestamptz{Time: getStartTime(input.Start), Valid: true}
	referrers, err := s.Repo.GetTopReferrers(ctx, db.GetTopReferrersParams{
		ShortUrlID: shortURL.ID,
		ClickedAt:  clickedAtParam,
		Limit:      5,
	})
	if err != nil {
		s.Logger.Err(err).Msg("failed to get referrers")
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

	return GetSummaryOutput{
		TotalClicks:  summary.TotalClicks,
		UniqueClicks: summary.UniqueClicks,
		BotClicks:    summary.BotClicks,
		TopReferrers: topReferrers,
	}, nil
}
