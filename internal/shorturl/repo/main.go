package repo

import (
	"context"

	db "url_shortner_backend/db/output"
)

type ShortURLRepository interface {
	CreateShortURL(ctx context.Context, params db.CreateShortURLParams) (db.CreateShortURLRow, error)
	UpdateShortURLCode(ctx context.Context, params db.UpdateShortURLCodeParams) (db.UpdateShortURLCodeRow, error)
	GetShortURLByCode(ctx context.Context, code string) (db.ShortUrl, error)
	GetShortURLsByOwner(ctx context.Context, params db.GetShortURLsByOwnerParams) ([]db.ShortUrl, error)
	TransferAnonymousURLsToUser(ctx context.Context, params db.TransferAnonymousURLsToUserParams) error
	CountURLsCreatedToday(ctx context.Context, ownerID string) (int64, error)
}

type ShortURLRepoImp struct {
	Queries *db.Queries
}

type ShortURLRepoParams struct {
	Queries *db.Queries
}

func NewShortURLRepoImp(p ShortURLRepoParams) *ShortURLRepoImp {
	return &ShortURLRepoImp{
		Queries: p.Queries,
	}
}

func (r *ShortURLRepoImp) CreateShortURL(ctx context.Context, params db.CreateShortURLParams) (db.CreateShortURLRow, error) {
	return r.Queries.CreateShortURL(ctx, params)
}

func (r *ShortURLRepoImp) UpdateShortURLCode(ctx context.Context, params db.UpdateShortURLCodeParams) (db.UpdateShortURLCodeRow, error) {
	return r.Queries.UpdateShortURLCode(ctx, params)
}

func (r *ShortURLRepoImp) GetShortURLByCode(ctx context.Context, code string) (db.ShortUrl, error) {
	return r.Queries.GetShortURLByCode(ctx, code)
}

func (r *ShortURLRepoImp) GetShortURLsByOwner(ctx context.Context, params db.GetShortURLsByOwnerParams) ([]db.ShortUrl, error) {
	return r.Queries.GetShortURLsByOwner(ctx, params)
}

func (r *ShortURLRepoImp) TransferAnonymousURLsToUser(ctx context.Context, params db.TransferAnonymousURLsToUserParams) error {
	return r.Queries.TransferAnonymousURLsToUser(ctx, params)
}

func (r *ShortURLRepoImp) CountURLsCreatedToday(ctx context.Context, ownerID string) (int64, error) {
	return r.Queries.CountURLsCreatedToday(ctx, ownerID)
}
