package service

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"

	db "url_shortner_backend/db/output"
	"url_shortner_backend/pkg/hash"
)

func (s *AuthSvcImp) Login(ctx context.Context, input LoginInput) (AuthOutput, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))

	user, err := s.Repo.GetUserByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return AuthOutput{}, ErrInvalidCredentials
		}
		s.Logger.Error("failed to get user", "email", email, "error", err)
		return AuthOutput{}, ErrInvalidCredentials
	}

	if !hash.CheckPassword(input.Password, user.PasswordHash) {
		return AuthOutput{}, ErrInvalidCredentials
	}

	userID := uuidToString(user.ID)

	// Transfer anonymous URLs if anonID provided
	if input.AnonID != "" {
		err = s.ShortURLRepo.TransferAnonymousURLsToUser(ctx, db.TransferAnonymousURLsToUserParams{
			OwnerID:   input.AnonID,
			OwnerID_2: userID,
		})
		if err != nil {
			s.Logger.Error("failed to transfer anonymous urls on login", "anon_id", input.AnonID, "user_id", userID, "error", err)
		} else {
			s.Logger.Info("transferred anonymous urls on login", "anon_id", input.AnonID, "user_id", userID)
		}
	}

	return s.generateTokens(ctx, userID, email)
}
