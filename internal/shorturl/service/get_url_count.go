package service

import (
	"context"
)

func (s *ShortURLSvcImp) GetUserURLCountThisMonth(ctx context.Context, userID string) (int64, error) {
	count, err := s.Repo.CountURLsCreatedThisMonth(ctx, userID)
	if err != nil {
		s.Logger.Error("failed to count urls this month", "user_id", userID, "error", err)
		return 0, ErrURLFetch
	}
	return count, nil
}
