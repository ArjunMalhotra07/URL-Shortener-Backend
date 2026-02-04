package service

import (
	"context"
	"time"
	db "url_shortner_backend/db/output"
)

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
