package service

import "url_shortner_backend/internal/shorturl/repo"

type ShortURLSvc interface{}

type ShortURLSvcImp struct {
	Repo repo.ShortURLRepository
}
