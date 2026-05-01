package repo

import (
	"context"

	db "url_shortner_backend/db/output"
)

type AdminRepository interface {
	ListUsers(ctx context.Context, params db.AdminListUsersParams) ([]db.AdminListUsersRow, error)
	CountUsers(ctx context.Context) (int64, error)
	GetUserURLs(ctx context.Context, params db.AdminGetUserURLsParams) ([]db.AdminGetUserURLsRow, error)
	CountUserURLs(ctx context.Context, ownerID string) (int64, error)
	GetPlatformStats(ctx context.Context) (db.AdminGetPlatformStatsRow, error)
	GetUsersByTier(ctx context.Context) ([]db.AdminGetUsersByTierRow, error)
}

type AdminRepoImp struct {
	Queries *db.Queries
}

type AdminRepoParams struct {
	Queries *db.Queries
}

func NewAdminRepoImp(p AdminRepoParams) *AdminRepoImp {
	return &AdminRepoImp{
		Queries: p.Queries,
	}
}

func (r *AdminRepoImp) ListUsers(ctx context.Context, params db.AdminListUsersParams) ([]db.AdminListUsersRow, error) {
	return r.Queries.AdminListUsers(ctx, params)
}

func (r *AdminRepoImp) CountUsers(ctx context.Context) (int64, error) {
	return r.Queries.AdminCountUsers(ctx)
}

func (r *AdminRepoImp) GetUserURLs(ctx context.Context, params db.AdminGetUserURLsParams) ([]db.AdminGetUserURLsRow, error) {
	return r.Queries.AdminGetUserURLs(ctx, params)
}

func (r *AdminRepoImp) CountUserURLs(ctx context.Context, ownerID string) (int64, error) {
	return r.Queries.AdminCountUserURLs(ctx, ownerID)
}

func (r *AdminRepoImp) GetPlatformStats(ctx context.Context) (db.AdminGetPlatformStatsRow, error) {
	return r.Queries.AdminGetPlatformStats(ctx)
}

func (r *AdminRepoImp) GetUsersByTier(ctx context.Context) ([]db.AdminGetUsersByTierRow, error) {
	return r.Queries.AdminGetUsersByTier(ctx)
}
