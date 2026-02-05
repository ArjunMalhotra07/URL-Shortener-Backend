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
	TransferAnonymousURLsToUserWithLimit(ctx context.Context, params db.TransferAnonymousURLsToUserWithLimitParams) error
	CountURLsCreatedThisMonth(ctx context.Context, ownerID string) (int64, error)
	CountURLsByOwner(ctx context.Context, params db.CountURLsByOwnerParams) (int64, error)
	GetURLByCodeAndOwner(ctx context.Context, params db.GetURLByCodeAndOwnerParams) (db.ShortUrl, error)
	ToggleURLActive(ctx context.Context, params db.ToggleURLActiveParams) error
	SoftDeleteURL(ctx context.Context, params db.SoftDeleteURLParams) error
	UpdateLongURL(ctx context.Context, params db.UpdateLongURLParams) (db.ShortUrl, error)
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

func (r *ShortURLRepoImp) TransferAnonymousURLsToUserWithLimit(ctx context.Context, params db.TransferAnonymousURLsToUserWithLimitParams) error {
	return r.Queries.TransferAnonymousURLsToUserWithLimit(ctx, params)
}

func (r *ShortURLRepoImp) CountURLsCreatedThisMonth(ctx context.Context, ownerID string) (int64, error) {
	return r.Queries.CountURLsCreatedThisMonth(ctx, ownerID)
}

func (r *ShortURLRepoImp) CountURLsByOwner(ctx context.Context, params db.CountURLsByOwnerParams) (int64, error) {
	return r.Queries.CountURLsByOwner(ctx, params)
}

func (r *ShortURLRepoImp) GetURLByCodeAndOwner(ctx context.Context, params db.GetURLByCodeAndOwnerParams) (db.ShortUrl, error) {
	return r.Queries.GetURLByCodeAndOwner(ctx, params)
}

func (r *ShortURLRepoImp) ToggleURLActive(ctx context.Context, params db.ToggleURLActiveParams) error {
	return r.Queries.ToggleURLActive(ctx, params)
}

func (r *ShortURLRepoImp) SoftDeleteURL(ctx context.Context, params db.SoftDeleteURLParams) error {
	return r.Queries.SoftDeleteURL(ctx, params)
}

func (r *ShortURLRepoImp) UpdateLongURL(ctx context.Context, params db.UpdateLongURLParams) (db.ShortUrl, error) {
	return r.Queries.UpdateLongURL(ctx, params)
}
