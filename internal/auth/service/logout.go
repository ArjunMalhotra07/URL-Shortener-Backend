package service

import (
	"context"
)

func (s *AuthSvcImp) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := hashToken(refreshToken)
	return s.Repo.DeleteRefreshToken(ctx, tokenHash)
}

func (s *AuthSvcImp) LogoutAll(ctx context.Context, userID string) error {
	uid := stringToUUID(userID)
	return s.Repo.DeleteAllUserRefreshTokens(ctx, uid)
}
