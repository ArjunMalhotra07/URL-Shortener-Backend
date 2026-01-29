package service

import (
	"context"
	"errors"
	"net"
	"net/url"
	"strings"

	"url_shortner_backend/internal/shorturl/repo"
	"url_shortner_backend/pkg/config"
	"url_shortner_backend/pkg/logger"
)

type ShortURLSvc interface {
	CreateShortURL(ctx context.Context, input CreateShortURLInput) (CreateShortURLOutput, error)
	GetLongURL(ctx context.Context, input GetLongURLInput) (GetLongURLOutput, error)
	GetMyURLs(ctx context.Context, input GetMyURLsInput) (GetMyURLsOutput, error)
	TransferURLsToUser(ctx context.Context, input TransferURLsInput) error
	ToggleURLActive(ctx context.Context, input ToggleURLInput) error
	DeleteURL(ctx context.Context, input DeleteURLInput) error
}

type ShortURLSvcImp struct {
	Logger logger.Logger
	Repo   repo.ShortURLRepository
	Cfg    *config.Config
}

func NewShortURLSvcImp(Repo repo.ShortURLRepository, Logger logger.Logger, cfg *config.Config) *ShortURLSvcImp {
	return &ShortURLSvcImp{
		Repo:   Repo,
		Logger: Logger,
		Cfg:    cfg,
	}
}

func normalizeURL(rawURL string) string {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		return "https://" + rawURL
	}
	return rawURL
}

const (
	base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	// Offset to ensure minimum 3-character codes (62^2 = 3844)
	base62Offset = 3844
)

func encodeBase62(num int64) string {
	num += base62Offset // Ensures 3+ character codes

	var encoded []byte
	for num > 0 {
		encoded = append([]byte{base62Chars[num%62]}, encoded...)
		num /= 62
	}
	return string(encoded)
}

func validateURL(raw string) error {
	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return errors.New("invalid URL format")
	}

	// 1️⃣ Allow only http/https
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("only http and https URLs are allowed")
	}

	// 2️⃣ Host must exist and look like a domain
	host := u.Hostname()
	if host == "" || !strings.Contains(host, ".") {
		return errors.New("invalid domain")
	}

	// 3️⃣ Block localhost & private IPs (security)
	ip := net.ParseIP(host)
	if ip != nil {
		if ip.IsLoopback() || ip.IsPrivate() {
			return errors.New("private or local addresses are not allowed")
		}
	}

	return nil
}
