package service

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pquerna/otp/totp"
	"github.com/rs/zerolog"

	"url_shortner_backend/internal/admin/repo"
	"url_shortner_backend/pkg/config"
	"url_shortner_backend/pkg/jwt"
)

type AdminService interface {
	Login(ctx context.Context, input LoginInput) (LoginOutput, error)
	ListUsers(ctx context.Context, input ListUsersInput) (ListUsersOutput, error)
	GetUserURLs(ctx context.Context, input GetUserURLsInput) (GetUserURLsOutput, error)
	GetPlatformStats(ctx context.Context) (PlatformStatsOutput, error)
}

type AdminSvcImp struct {
	Repo   repo.AdminRepository
	JWT    *jwt.JWTManager
	Cfg    *config.Config
	Logger zerolog.Logger
}

func NewAdminSvcImp(r repo.AdminRepository, jwtMgr *jwt.JWTManager, cfg *config.Config, l zerolog.Logger) *AdminSvcImp {
	return &AdminSvcImp{
		Repo:   r,
		JWT:    jwtMgr,
		Cfg:    cfg,
		Logger: l,
	}
}

// --- Login ---

type LoginInput struct {
	AdminID string
	Code    string
}

type LoginOutput struct {
	AccessToken     string
	AccessExpiresAt time.Time
}

func (s *AdminSvcImp) Login(ctx context.Context, input LoginInput) (LoginOutput, error) {
	if input.AdminID != s.Cfg.AdminID {
		return LoginOutput{}, ErrInvalidCredentials
	}

	valid := totp.Validate(input.Code, s.Cfg.AdminTOTPSecret)
	if !valid {
		return LoginOutput{}, ErrInvalidCredentials
	}

	// Generate an access token with a special admin user ID
	accessToken, expiresAt, err := s.JWT.GenerateAccessToken("admin", "admin@tinyclk")
	if err != nil {
		s.Logger.Err(err).Msg("failed to generate admin access token")
		return LoginOutput{}, ErrTokenCreation
	}

	s.Logger.Info().Msg("admin logged in")

	return LoginOutput{
		AccessToken:     accessToken,
		AccessExpiresAt: expiresAt,
	}, nil
}

// --- List Users ---

type ListUsersInput struct {
	Limit  int32
	Offset int32
}

type UserRow struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Tier          string `json:"tier"`
	LoginType     int16  `json:"login_type"`
	CreatedAt     string `json:"created_at"`
	UrlCount      int64  `json:"url_count"`
	ActiveCount   int64  `json:"active_count"`
	InactiveCount int64  `json:"inactive_count"`
	DeletedCount  int64  `json:"deleted_count"`
	TotalClicks   int64  `json:"total_clicks"`
}

type ListUsersOutput struct {
	Users []UserRow `json:"users"`
	Total int64     `json:"total"`
}

func (s *AdminSvcImp) ListUsers(ctx context.Context, input ListUsersInput) (ListUsersOutput, error) {
	if input.Limit <= 0 {
		input.Limit = 20
	}
	if input.Limit > 100 {
		input.Limit = 100
	}

	users, err := s.Repo.ListUsers(ctx, struct {
		Limit  int32 `json:"limit"`
		Offset int32 `json:"offset"`
	}{Limit: input.Limit, Offset: input.Offset})
	if err != nil {
		s.Logger.Err(err).Msg("failed to list users")
		return ListUsersOutput{}, err
	}

	total, err := s.Repo.CountUsers(ctx)
	if err != nil {
		s.Logger.Err(err).Msg("failed to count users")
		return ListUsersOutput{}, err
	}

	rows := make([]UserRow, len(users))
	for i, u := range users {
		rows[i] = UserRow{
			ID:            uuidToString(u.ID),
			Email:         u.Email,
			Name:          textToString(u.Name),
			Tier:          string(u.Tier),
			LoginType:     u.LoginType,
			CreatedAt:     u.CreatedAt.Time.UTC().Format(time.RFC3339),
			UrlCount:      u.UrlCount,
			ActiveCount:   u.ActiveCount,
			InactiveCount: u.InactiveCount,
			DeletedCount:  u.DeletedCount,
			TotalClicks:   u.TotalClicks,
		}
	}

	return ListUsersOutput{Users: rows, Total: total}, nil
}

// --- Get User URLs ---

type GetUserURLsInput struct {
	UserID string
	Limit  int32
	Offset int32
}

type URLRow struct {
	ID         int64  `json:"id"`
	Code       string `json:"code"`
	LongUrl    string `json:"long_url"`
	Name       string `json:"name"`
	IsActive   bool   `json:"is_active"`
	IsDeleted  bool   `json:"is_deleted"`
	CreatedAt  string `json:"created_at"`
	ExpiresAt  string `json:"expires_at,omitempty"`
	ClickCount int64  `json:"click_count"`
}

type GetUserURLsOutput struct {
	URLs  []URLRow `json:"urls"`
	Total int64    `json:"total"`
}

func (s *AdminSvcImp) GetUserURLs(ctx context.Context, input GetUserURLsInput) (GetUserURLsOutput, error) {
	if input.Limit <= 0 {
		input.Limit = 20
	}
	if input.Limit > 100 {
		input.Limit = 100
	}

	urls, err := s.Repo.GetUserURLs(ctx, struct {
		OwnerID string `json:"owner_id"`
		Limit   int32  `json:"limit"`
		Offset  int32  `json:"offset"`
	}{OwnerID: input.UserID, Limit: input.Limit, Offset: input.Offset})
	if err != nil {
		s.Logger.Err(err).Msg("failed to get user urls")
		return GetUserURLsOutput{}, err
	}

	total, err := s.Repo.CountUserURLs(ctx, input.UserID)
	if err != nil {
		s.Logger.Err(err).Msg("failed to count user urls")
		return GetUserURLsOutput{}, err
	}

	rows := make([]URLRow, len(urls))
	for i, u := range urls {
		row := URLRow{
			ID:         u.ID,
			Code:       u.Code,
			LongUrl:    u.LongUrl,
			Name:       textToString(u.Name),
			IsActive:   u.IsActive,
			IsDeleted:  u.IsDeleted,
			CreatedAt:  u.CreatedAt.Time.UTC().Format(time.RFC3339),
			ClickCount: u.ClickCount,
		}
		if u.ExpiresAt.Valid {
			row.ExpiresAt = u.ExpiresAt.Time.UTC().Format(time.RFC3339)
		}
		rows[i] = row
	}

	return GetUserURLsOutput{URLs: rows, Total: total}, nil
}

// --- Platform Stats ---

type TierCount struct {
	Tier  string `json:"tier"`
	Count int64  `json:"count"`
}

type PlatformStatsOutput struct {
	TotalUsers   int64       `json:"total_users"`
	TotalUrls    int64       `json:"total_urls"`
	ActiveUrls   int64       `json:"active_urls"`
	InactiveUrls int64       `json:"inactive_urls"`
	DeletedUrls  int64       `json:"deleted_urls"`
	TotalClicks  int64       `json:"total_clicks"`
	UsersByTier  []TierCount `json:"users_by_tier"`
}

func (s *AdminSvcImp) GetPlatformStats(ctx context.Context) (PlatformStatsOutput, error) {
	stats, err := s.Repo.GetPlatformStats(ctx)
	if err != nil {
		s.Logger.Err(err).Msg("failed to get platform stats")
		return PlatformStatsOutput{}, err
	}

	tiers, err := s.Repo.GetUsersByTier(ctx)
	if err != nil {
		s.Logger.Err(err).Msg("failed to get users by tier")
		return PlatformStatsOutput{}, err
	}

	tierCounts := make([]TierCount, len(tiers))
	for i, t := range tiers {
		tierCounts[i] = TierCount{
			Tier:  string(t.Tier),
			Count: t.Count,
		}
	}

	return PlatformStatsOutput{
		TotalUsers:   stats.TotalUsers,
		TotalUrls:    stats.TotalUrls,
		ActiveUrls:   stats.ActiveUrls,
		InactiveUrls: stats.InactiveUrls,
		DeletedUrls:  stats.DeletedUrls,
		TotalClicks:  stats.TotalClicks,
		UsersByTier:  tierCounts,
	}, nil
}

// Helpers

func uuidToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	return hex.EncodeToString(u.Bytes[0:4]) + "-" +
		hex.EncodeToString(u.Bytes[4:6]) + "-" +
		hex.EncodeToString(u.Bytes[6:8]) + "-" +
		hex.EncodeToString(u.Bytes[8:10]) + "-" +
		hex.EncodeToString(u.Bytes[10:16])
}

func textToString(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}
