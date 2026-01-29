package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	db "url_shortner_backend/db/output"
	"url_shortner_backend/internal/auth/repo"
	shorturlrepo "url_shortner_backend/internal/shorturl/repo"
	"url_shortner_backend/pkg/config"
	"url_shortner_backend/pkg/jwt"
	"url_shortner_backend/pkg/logger"
)

type AuthService interface {
	Signup(ctx context.Context, input SignupInput) (AuthOutput, error)
	Login(ctx context.Context, input LoginInput) (AuthOutput, error)
	Refresh(ctx context.Context, refreshToken string) (AuthOutput, error)
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, userID string) error
}

type AuthSvcImp struct {
	Repo         repo.AuthRepository
	ShortURLRepo shorturlrepo.ShortURLRepository
	JWT          *jwt.JWTManager
	Logger       logger.Logger
	Cfg          *config.Config
}

func NewAuthSvcImp(r repo.AuthRepository, shortURLRepo shorturlrepo.ShortURLRepository, jwtMgr *jwt.JWTManager, l logger.Logger, cfg *config.Config) *AuthSvcImp {
	return &AuthSvcImp{
		Repo:         r,
		ShortURLRepo: shortURLRepo,
		JWT:          jwtMgr,
		Logger:       l,
		Cfg:          cfg,
	}
}

type SignupInput struct {
	Email    string
	Password string
	AnonID   string
}

type LoginInput struct {
	Email    string
	Password string
	AnonID   string
}

type AuthOutput struct {
	UserID           string
	Email            string
	AccessToken      string
	AccessExpiresAt  time.Time
	RefreshToken     string
	RefreshExpiresAt time.Time
}

func (s *AuthSvcImp) generateTokens(ctx context.Context, userID, email string) (AuthOutput, error) {
	accessToken, accessExpiresAt, err := s.JWT.GenerateAccessToken(userID, email)
	if err != nil {
		s.Logger.Error("failed to generate access token", "error", err)
		return AuthOutput{}, ErrTokenCreation
	}

	refreshToken, refreshExpiresAt, err := s.JWT.GenerateRefreshToken()
	if err != nil {
		s.Logger.Error("failed to generate refresh token", "error", err)
		return AuthOutput{}, ErrTokenCreation
	}

	tokenHash := hashToken(refreshToken)
	_, err = s.Repo.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		UserID:    stringToUUID(userID),
		TokenHash: tokenHash,
		ExpiresAt: pgtype.Timestamptz{Time: refreshExpiresAt, Valid: true},
	})
	if err != nil {
		s.Logger.Error("failed to store refresh token", "error", err)
		return AuthOutput{}, ErrTokenCreation
	}

	s.Logger.Info("tokens generated", "user_id", userID)

	return AuthOutput{
		UserID:           userID,
		Email:            email,
		AccessToken:      accessToken,
		AccessExpiresAt:  accessExpiresAt,
		RefreshToken:     refreshToken,
		RefreshExpiresAt: refreshExpiresAt,
	}, nil
}

// transferAnonymousURLsWithQuota transfers anonymous URLs to user, respecting monthly quota.
// Only transfers up to (quota - existing_this_month) URLs. Excess URLs are not transferred.
func (s *AuthSvcImp) transferAnonymousURLsWithQuota(ctx context.Context, anonID, userID string) {
	// Count how many URLs user has already created this month
	monthCount, err := s.ShortURLRepo.CountURLsCreatedThisMonth(ctx, userID)
	if err != nil {
		s.Logger.Error("failed to count user urls this month", "user_id", userID, "error", err)
		return
	}

	remaining := s.Cfg.MonthlyQuotaUser - int(monthCount)
	if remaining <= 0 {
		s.Logger.Info("user at monthly quota, skipping anonymous url transfer", "user_id", userID, "month_count", monthCount)
		return
	}

	// Transfer only up to remaining quota
	err = s.ShortURLRepo.TransferAnonymousURLsToUserWithLimit(ctx, db.TransferAnonymousURLsToUserWithLimitParams{
		OwnerID:   anonID,
		OwnerID_2: userID,
		Limit:     int32(remaining),
	})
	if err != nil {
		s.Logger.Error("failed to transfer anonymous urls", "anon_id", anonID, "user_id", userID, "error", err)
		return
	}

	s.Logger.Info("transferred anonymous urls with quota limit", "anon_id", anonID, "user_id", userID, "max_transferred", remaining)
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

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

func stringToUUID(s string) pgtype.UUID {
	s = strings.ReplaceAll(s, "-", "")
	bytes, _ := hex.DecodeString(s)
	var uuid pgtype.UUID
	if len(bytes) == 16 {
		copy(uuid.Bytes[:], bytes)
		uuid.Valid = true
	}
	return uuid
}
