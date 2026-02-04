package repo

import (
	"context"

	db "url_shortner_backend/db/output"
)

type AnalyticsRepository interface {
	InsertClick(ctx context.Context, params db.InsertClickParams) error
	GetClicksByShortURLID(ctx context.Context, params db.GetClicksByShortURLIDParams) ([]db.Click, error)
	CountClicksByShortURLID(ctx context.Context, shortUrlID int64) (int64, error)
	CountUniqueClicksByShortURLID(ctx context.Context, shortUrlID int64) (int64, error)
	GetClicksByCountry(ctx context.Context, params db.GetClicksByCountryParams) ([]db.GetClicksByCountryRow, error)
	GetClicksByCity(ctx context.Context, params db.GetClicksByCityParams) ([]db.GetClicksByCityRow, error)
	GetClicksByDevice(ctx context.Context, params db.GetClicksByDeviceParams) ([]db.GetClicksByDeviceRow, error)
	GetClicksByBrowser(ctx context.Context, params db.GetClicksByBrowserParams) ([]db.GetClicksByBrowserRow, error)
	GetClicksByOS(ctx context.Context, params db.GetClicksByOSParams) ([]db.GetClicksByOSRow, error)
	GetClicksTimeseries(ctx context.Context, params db.GetClicksTimeseriesParams) ([]db.GetClicksTimeseriesRow, error)
	GetClicksTimeseriesHourly(ctx context.Context, params db.GetClicksTimeseriesHourlyParams) ([]db.GetClicksTimeseriesHourlyRow, error)
	GetTopReferrers(ctx context.Context, params db.GetTopReferrersParams) ([]db.GetTopReferrersRow, error)
	GetTopCampaigns(ctx context.Context, params db.GetTopCampaignsParams) ([]db.GetTopCampaignsRow, error)
	GetClicksSummary(ctx context.Context, params db.GetClicksSummaryParams) (db.GetClicksSummaryRow, error)
	GetClicksSummaryAllTime(ctx context.Context, shortUrlID int64) (db.GetClicksSummaryAllTimeRow, error)
}

type AnalyticsRepoImp struct {
	Queries *db.Queries
}

type AnalyticsRepoParams struct {
	Queries *db.Queries
}

func NewAnalyticsRepoImp(p AnalyticsRepoParams) *AnalyticsRepoImp {
	return &AnalyticsRepoImp{
		Queries: p.Queries,
	}
}

func (r *AnalyticsRepoImp) InsertClick(ctx context.Context, params db.InsertClickParams) error {
	return r.Queries.InsertClick(ctx, params)
}

func (r *AnalyticsRepoImp) GetClicksByShortURLID(ctx context.Context, params db.GetClicksByShortURLIDParams) ([]db.Click, error) {
	return r.Queries.GetClicksByShortURLID(ctx, params)
}

func (r *AnalyticsRepoImp) CountClicksByShortURLID(ctx context.Context, shortUrlID int64) (int64, error) {
	return r.Queries.CountClicksByShortURLID(ctx, shortUrlID)
}

func (r *AnalyticsRepoImp) CountUniqueClicksByShortURLID(ctx context.Context, shortUrlID int64) (int64, error) {
	return r.Queries.CountUniqueClicksByShortURLID(ctx, shortUrlID)
}

func (r *AnalyticsRepoImp) GetClicksByCountry(ctx context.Context, params db.GetClicksByCountryParams) ([]db.GetClicksByCountryRow, error) {
	return r.Queries.GetClicksByCountry(ctx, params)
}

func (r *AnalyticsRepoImp) GetClicksByCity(ctx context.Context, params db.GetClicksByCityParams) ([]db.GetClicksByCityRow, error) {
	return r.Queries.GetClicksByCity(ctx, params)
}

func (r *AnalyticsRepoImp) GetClicksByDevice(ctx context.Context, params db.GetClicksByDeviceParams) ([]db.GetClicksByDeviceRow, error) {
	return r.Queries.GetClicksByDevice(ctx, params)
}

func (r *AnalyticsRepoImp) GetClicksByBrowser(ctx context.Context, params db.GetClicksByBrowserParams) ([]db.GetClicksByBrowserRow, error) {
	return r.Queries.GetClicksByBrowser(ctx, params)
}

func (r *AnalyticsRepoImp) GetClicksByOS(ctx context.Context, params db.GetClicksByOSParams) ([]db.GetClicksByOSRow, error) {
	return r.Queries.GetClicksByOS(ctx, params)
}

func (r *AnalyticsRepoImp) GetClicksTimeseries(ctx context.Context, params db.GetClicksTimeseriesParams) ([]db.GetClicksTimeseriesRow, error) {
	return r.Queries.GetClicksTimeseries(ctx, params)
}

func (r *AnalyticsRepoImp) GetClicksTimeseriesHourly(ctx context.Context, params db.GetClicksTimeseriesHourlyParams) ([]db.GetClicksTimeseriesHourlyRow, error) {
	return r.Queries.GetClicksTimeseriesHourly(ctx, params)
}

func (r *AnalyticsRepoImp) GetTopReferrers(ctx context.Context, params db.GetTopReferrersParams) ([]db.GetTopReferrersRow, error) {
	return r.Queries.GetTopReferrers(ctx, params)
}

func (r *AnalyticsRepoImp) GetTopCampaigns(ctx context.Context, params db.GetTopCampaignsParams) ([]db.GetTopCampaignsRow, error) {
	return r.Queries.GetTopCampaigns(ctx, params)
}

func (r *AnalyticsRepoImp) GetClicksSummary(ctx context.Context, params db.GetClicksSummaryParams) (db.GetClicksSummaryRow, error) {
	return r.Queries.GetClicksSummary(ctx, params)
}

func (r *AnalyticsRepoImp) GetClicksSummaryAllTime(ctx context.Context, shortUrlID int64) (db.GetClicksSummaryAllTimeRow, error) {
	return r.Queries.GetClicksSummaryAllTime(ctx, shortUrlID)
}
