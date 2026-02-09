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
			s.Logger.Info().Str("email", email).Msg("login attempt for non-existent user")
			return AuthOutput{}, ErrInvalidCredentials
		}
		s.Logger.Err(err).Str("email", email).Msg("failed to get user")
		return AuthOutput{}, ErrInvalidCredentials
	}

	// Check if user signed up with Google
	if user.LoginType == 1 {
		s.Logger.Info().Str("email", email).Msg("password login attempt for google user")
		return AuthOutput{}, ErrEmailExistsWithGoogle
	}

	if !hash.CheckPassword(input.Password, user.PasswordHash.String) {
		s.Logger.Info().Str("email", email).Msg("login attempt with wrong password")
		return AuthOutput{}, ErrInvalidCredentials
	}

	userID := uuidToString(user.ID)

	// Transfer anonymous URLs if anonID provided (with quota limit)
	if input.AnonID != "" {
		s.ShortURLSvc.TransferAnonymousURLsWithQuota(ctx, input.AnonID, userID)
	}

	s.Logger.Info().Str("email", email).Str("user_id", userID).Msg("user logged in successfully")
	return s.generateTokens(ctx, userID, email)
}
