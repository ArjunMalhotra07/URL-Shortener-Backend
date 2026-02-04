package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"

	db "url_shortner_backend/db/output"
	"url_shortner_backend/pkg/useragent"

	"github.com/jackc/pgx/v5/pgtype"
)

type RecordClickInput struct {
	ShortURLID  int64
	IPAddress   string
	UserAgent   string
	Referrer    string
	UTMSource   string
	UTMMedium   string
	UTMCampaign string
}

func (s *AnalyticsSvcImp) RecordClick(ctx context.Context, input RecordClickInput) error {
	// Parse user agent
	parsedUA := useragent.Parse(input.UserAgent)
	// ipAddress := "8.8.8.8"
	ipAddress := input.IPAddress
	// Lookup geo info
	var country, city, region string
	if s.GeoIP != nil && ipAddress != "" {
		geoResult, err := s.GeoIP.Lookup(ipAddress)
		if err == nil && geoResult != nil {
			country = geoResult.Country
			city = geoResult.City
			region = geoResult.Region
		}
	}

	// Hash IP for privacy
	ipHash := hashIP(ipAddress, input.ShortURLID)

	// Extract referrer domain
	referrerDomain := extractDomain(input.Referrer)

	// Check if unique visitor (using Redis)
	isUnique := false
	if s.Redis != nil {
		uniqueKey := fmt.Sprintf("unique:%d:%s", input.ShortURLID, ipHash)
		wasSet, err := s.Redis.SetNX(ctx, uniqueKey, "1", 24*time.Hour)
		if err == nil && wasSet {
			isUnique = true
		}

		// Invalidate summary cache
		cacheKey := fmt.Sprintf("analytics:%d:summary", input.ShortURLID)
		_ = s.Redis.Del(ctx, cacheKey)
	} else {
		// Without Redis, we can't efficiently track uniqueness
		// Mark as unique by default
		isUnique = true
	}

	// Insert click record
	params := db.InsertClickParams{
		ShortUrlID:     input.ShortURLID,
		IpHash:         ipHash,
		Country:        toPgText(country),
		City:           toPgText(city),
		Region:         toPgText(region),
		Browser:        toPgText(parsedUA.Browser),
		Os:             toPgText(parsedUA.OS),
		DeviceType:     toPgText(parsedUA.DeviceType),
		Referrer:       toPgText(input.Referrer),
		ReferrerDomain: toPgText(referrerDomain),
		UtmSource:      toPgText(input.UTMSource),
		UtmMedium:      toPgText(input.UTMMedium),
		UtmCampaign:    toPgText(input.UTMCampaign),
		IsUnique:       pgtype.Bool{Bool: isUnique, Valid: true},
		IsBot:          pgtype.Bool{Bool: parsedUA.IsBot, Valid: true},
	}

	err := s.Repo.InsertClick(ctx, params)
	if err != nil {
		s.Logger.Error("failed to record click", "error", err, "short_url_id", input.ShortURLID)
		return ErrClickRecording
	}

	s.Logger.Debug("click recorded", "short_url_id", input.ShortURLID, "country", country, "device", parsedUA.DeviceType)
	return nil
}

func hashIP(ip string, urlID int64) string {
	// Use URL ID as salt to prevent cross-URL tracking
	data := fmt.Sprintf("%s:%d", ip, urlID)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16]) // Use first 16 bytes
}

func extractDomain(referrer string) string {
	if referrer == "" {
		return ""
	}
	u, err := url.Parse(referrer)
	if err != nil {
		return ""
	}
	host := u.Hostname()
	// Remove www. prefix
	host = strings.TrimPrefix(host, "www.")
	return host
}

func toPgText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}
