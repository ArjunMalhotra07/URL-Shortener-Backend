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
		s.Logger.Error("failed to check existing user", "email", email, "error", err)
		return AuthOutput{}, ErrUserCreation
	}

	// Hash password
	passwordHash, err := hash.HashPassword(input.Password)
	if err != nil {
		s.Logger.Error("failed to hash password", "error", err)
		return AuthOutput{}, ErrUserCreation
	}

	// Create user
	user, err := s.Repo.CreateUser(ctx, db.CreateUserParams{
		Email:        email,
		PasswordHash: pgtype.Text{String: passwordHash, Valid: true},
	})
	if err != nil {
		s.Logger.Error("failed to create user", "email", email, "error", err)
		return AuthOutput{}, ErrUserCreation
	}

	userID := uuidToString(user.ID)

	// Transfer anonymous URLs if anonID provided (with quota limit)
	if input.AnonID != "" {
		s.transferAnonymousURLsWithQuota(ctx, input.AnonID, userID)
	}

	return s.generateTokens(ctx, userID, email)
}
