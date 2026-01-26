package repo

import db "url_shortner_backend/db/output"

type ShortURLRepository interface{}

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
