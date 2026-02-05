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
	shorturlsvc "url_shortner_backend/internal/shorturl/service"
	"url_shortner_backend/pkg/jwt"
	"url_shortner_backend/pkg/logger"
)

type AuthService interface {
	Signup(ctx context.Context, input SignupInput) (AuthOutput, error)
	Login(ctx context.Context, input LoginInput) (AuthOutput, error)
	GoogleLogin(ctx context.Context, input GoogleLoginInput) (AuthOutput, error)
	Refresh(ctx context.Context, refreshToken string) (AuthOutput, error)
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, userID string) error
	GetMe(ctx context.Context, userID string) (MeOutput, error)
}

type MeOutput struct {
	UserID              string
	Email               string
	Name                string
	AvatarURL           string
	Tier                string
	SubscriptionEndsAt  *time.Time
	URLsCreatedThisMonth int64
	URLsLimit           int
}

type AuthSvcImp struct {
	Repo        repo.AuthRepository
	ShortURLSvc shorturlsvc.ShortURLSvc
	JWT         *jwt.JWTManager
	Logger      logger.Logger
}

func NewAuthSvcImp(r repo.AuthRepository, shortURLSvc shorturlsvc.ShortURLSvc, jwtMgr *jwt.JWTManager, l logger.Logger) *AuthSvcImp {
	return &AuthSvcImp{
		Repo:        r,
		ShortURLSvc: shortURLSvc,
		JWT:         jwtMgr,
		Logger:      l,
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

type GoogleLoginInput struct {
	GoogleID  string
	Email     string
	Name      string
	AvatarURL string
	AnonID    string
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
