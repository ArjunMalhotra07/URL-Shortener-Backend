package service

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	db "url_shortner_backend/db/output"
)

func (s *AuthSvcImp) GoogleLogin(ctx context.Context, input GoogleLoginInput) (AuthOutput, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))

	user, err := s.Repo.GetUserByGoogleID(ctx, pgtype.Text{String: input.GoogleID, Valid: true})
	if err == nil {
		userID := uuidToString(user.ID)
		if input.AnonID != "" {
			s.ShortURLSvc.TransferAnonymousURLsWithQuota(ctx, input.AnonID, userID)
		}
		s.Logger.Info("google user logged in", "email", email, "user_id", userID)
		return s.generateTokens(ctx, userID, email)
	}

	if err != pgx.ErrNoRows {
		s.Logger.Error("failed to check google user", "google_id", input.GoogleID, "error", err)
		return AuthOutput{}, ErrUserCreation
	}

	existingUser, err := s.Repo.GetUserByEmail(ctx, email)
	if err == nil {
		if existingUser.LoginType == 0 {
			s.Logger.Info("google login attempt for email/password user", "email", email)
			return AuthOutput{}, ErrEmailExistsWithPassword
		}
		s.Logger.Error("email exists with different google account", "email", email)
		return AuthOutput{}, ErrEmailExists
	}

	if err != pgx.ErrNoRows {
		s.Logger.Error("failed to check existing user by email", "email", email, "error", err)
		return AuthOutput{}, ErrUserCreation
	}

	newUser, err := s.Repo.CreateGoogleUser(ctx, db.CreateGoogleUserParams{
		Email:     email,
		GoogleID:  pgtype.Text{String: input.GoogleID, Valid: true},
		Name:      pgtype.Text{String: input.Name, Valid: input.Name != ""},
		AvatarUrl: pgtype.Text{String: input.AvatarURL, Valid: input.AvatarURL != ""},
	})
	if err != nil {
		s.Logger.Error("failed to create google user", "email", email, "error", err)
		return AuthOutput{}, ErrUserCreation
	}

	userID := uuidToString(newUser.ID)

	if input.AnonID != "" {
		s.ShortURLSvc.TransferAnonymousURLsWithQuota(ctx, input.AnonID, userID)
	}

	s.Logger.Info("google user created and logged in", "email", email, "user_id", userID)
	return s.generateTokens(ctx, userID, email)
}
