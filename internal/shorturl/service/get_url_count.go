package service

import (
	"context"
)

func (s *ShortURLSvcImp) GetUserURLCountThisMonth(ctx context.Context, userID string) (int64, error) {
	count, err := s.Repo.CountURLsCreatedThisMonth(ctx, userID)
	if err != nil {
		s.Logger.Err(err).Str("user_id", userID).Msg("failed to count urls this month")
		return 0, ErrURLFetch
	}
	return count, nil
}
