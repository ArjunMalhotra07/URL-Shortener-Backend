package repo

import (
	"context"

	db "url_shortner_backend/db/output"
)

type ShortURLRepository interface {
	CreateShortURL(ctx context.Context, params db.CreateShortURLParams) (db.CreateShortURLRow, error)
	UpdateShortURLCode(ctx context.Context, params db.UpdateShortURLCodeParams) (db.UpdateShortURLCodeRow, error)
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
