package service

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"

	db "url_shortner_backend/db/output"
	"url_shortner_backend/pkg/hash"
)

func (s *AuthSvcImp) Signup(ctx context.Context, input SignupInput) (AuthOutput, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))

	// Check if user already exists
	_, err := s.Repo.GetUserByEmail(ctx, email)
	if err == nil {
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
		PasswordHash: passwordHash,
	})
	if err != nil {
		s.Logger.Error("failed to create user", "email", email, "error", err)
		return AuthOutput{}, ErrUserCreation
	}

	userID := uuidToString(user.ID)

	// Transfer anonymous URLs if anonID provided
	if input.AnonID != "" {
		err = s.ShortURLRepo.TransferAnonymousURLsToUser(ctx, db.TransferAnonymousURLsToUserParams{
			OwnerID:   input.AnonID,
			OwnerID_2: userID,
		})
		if err != nil {
			s.Logger.Error("failed to transfer anonymous urls", "anon_id", input.AnonID, "user_id", userID, "error", err)
			// Don't fail signup for this
		} else {
			s.Logger.Info("transferred anonymous urls", "anon_id", input.AnonID, "user_id", userID)
		}
	}

	return s.generateTokens(ctx, userID, email)
}
