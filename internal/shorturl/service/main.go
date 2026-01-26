package service

import (
	"context"
	"strings"

	"url_shortner_backend/internal/shorturl/repo"
	"url_shortner_backend/pkg/logger"
)

type ShortURLSvc interface {
	CreateShortURL(ctx context.Context, input CreateShortURLInput) (CreateShortURLOutput, error)
	GetLongURL(ctx context.Context, input GetLongURLInput) (GetLongURLOutput, error)
	GetMyURLs(ctx context.Context, input GetMyURLsInput) (GetMyURLsOutput, error)
	TransferURLsToUser(ctx context.Context, input TransferURLsInput) error
}

type ShortURLSvcImp struct {
	Logger logger.Logger
	Repo   repo.ShortURLRepository
}

func NewShortURLSvcImp(Repo repo.ShortURLRepository, Logger logger.Logger) *ShortURLSvcImp {
	return &ShortURLSvcImp{
		Repo:   Repo,
		Logger: Logger,
	}
}

func normalizeURL(rawURL string) string {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		return "https://" + rawURL
	}
	return rawURL
}

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func encodeBase62(num int64) string {
	if num == 0 {
		return string(base62Chars[0])
	}

	var encoded []byte
	for num > 0 {
		encoded = append([]byte{base62Chars[num%62]}, encoded...)
		num /= 62
	}
	return string(encoded)
}
