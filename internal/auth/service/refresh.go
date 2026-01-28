package service

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (s *AuthSvcImp) Refresh(ctx context.Context, refreshToken string) (AuthOutput, error) {
	tokenHash := hashToken(refreshToken)

	storedToken, err := s.Repo.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return AuthOutput{}, ErrInvalidToken
		}
		s.Logger.Error("failed to get refresh token", "error", err)
		return AuthOutput{}, ErrInvalidToken
	}

	// Delete old refresh token (rotation)
	_ = s.Repo.DeleteRefreshToken(ctx, tokenHash)

	user, err := s.Repo.GetUserByID(ctx, storedToken.UserID)
	if err != nil {
		s.Logger.Error("failed to get user for refresh", "error", err)
		return AuthOutput{}, ErrUserNotFound
	}

	userID := uuidToString(user.ID)

	return s.generateTokens(ctx, userID, user.Email)
}
