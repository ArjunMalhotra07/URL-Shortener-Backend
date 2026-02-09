package service

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	db "url_shortner_backend/db/output"
	"url_shortner_backend/pkg/hash"
)

func (s *AuthSvcImp) Signup(ctx context.Context, input SignupInput) (AuthOutput, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))

	// Check if user already exists
	existingUser, err := s.Repo.GetUserByEmail(ctx, email)
	if err == nil {
		if existingUser.LoginType == 1 {
			return AuthOutput{}, ErrEmailExistsWithGoogle
		}
		return AuthOutput{}, ErrEmailExists
	}
	if err != pgx.ErrNoRows {
		s.Logger.Err(err).Str("email", email).Msg("failed to check existing user")
		return AuthOutput{}, ErrUserCreation
	}

	// Hash password
	passwordHash, err := hash.HashPassword(input.Password)
	if err != nil {
		s.Logger.Err(err).Msg("failed to hash password")
		return AuthOutput{}, ErrUserCreation
	}

	// Create user
	user, err := s.Repo.CreateUser(ctx, db.CreateUserParams{
		Email:        email,
		PasswordHash: pgtype.Text{String: passwordHash, Valid: true},
	})
	if err != nil {
		s.Logger.Err(err).Str("email", email).Msg("failed to create user")
		return AuthOutput{}, ErrUserCreation
	}

	userID := uuidToString(user.ID)

	// Transfer anonymous URLs if anonID provided (with quota limit)
	if input.AnonID != "" {
		s.ShortURLSvc.TransferAnonymousURLsWithQuota(ctx, input.AnonID, userID)
	}

	return s.generateTokens(ctx, userID, email)
}
