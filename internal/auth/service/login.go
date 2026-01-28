package service

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"

	"url_shortner_backend/pkg/hash"
)

func (s *AuthSvcImp) Login(ctx context.Context, input LoginInput) (AuthOutput, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))

	user, err := s.Repo.GetUserByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			s.Logger.Info("login attempt for non-existent user", "email", email)
			return AuthOutput{}, ErrInvalidCredentials
		}
		s.Logger.Error("failed to get user", "email", email, "error", err)
		return AuthOutput{}, ErrInvalidCredentials
	}

	if !hash.CheckPassword(input.Password, user.PasswordHash) {
		s.Logger.Info("login attempt with wrong password", "email", email)
		return AuthOutput{}, ErrInvalidCredentials
	}

	userID := uuidToString(user.ID)

	// Transfer anonymous URLs if anonID provided (with quota limit)
	if input.AnonID != "" {
		s.transferAnonymousURLsWithQuota(ctx, input.AnonID, userID)
	}

	s.Logger.Info("user logged in successfully", "email", email, "user_id", userID)
	return s.generateTokens(ctx, userID, email)
}
