package service

import (
	"context"

	db "url_shortner_backend/db/output"
)

type TransferURLsInput struct {
	AnonID string
	UserID string
}

func (s *ShortURLSvcImp) TransferURLsToUser(ctx context.Context, input TransferURLsInput) error {
	if input.AnonID == "" || input.UserID == "" {
		s.Logger.Error("empty anon_id or user_id")
		return ErrInvalidOwner
	}

	err := s.Repo.TransferAnonymousURLsToUser(ctx, db.TransferAnonymousURLsToUserParams{
		OwnerID:   input.AnonID,
		OwnerID_2: input.UserID,
	})
	if err != nil {
		s.Logger.Error("failed to transfer urls", "anon_id", input.AnonID, "user_id", input.UserID, "error", err)
		return ErrURLTransfer
	}

	s.Logger.Info("urls transferred", "anon_id", input.AnonID, "user_id", input.UserID)
	return nil
}

// TransferAnonymousURLsWithQuota transfers anonymous URLs to user, respecting monthly quota.
// Also clears expiry date as a reward for signing up.
func (s *ShortURLSvcImp) TransferAnonymousURLsWithQuota(ctx context.Context, anonID, userID string) {
	if anonID == "" || userID == "" {
		return
	}

	// Count how many URLs user has already created this month
	monthCount, err := s.Repo.CountURLsCreatedThisMonth(ctx, userID)
	if err != nil {
		s.Logger.Error("failed to count user urls this month", "user_id", userID, "error", err)
		return
	}

	remaining := s.Cfg.MonthlyQuotaUser - int(monthCount)
	if remaining <= 0 {
		s.Logger.Info("user at monthly quota, skipping anonymous url transfer", "user_id", userID, "month_count", monthCount)
		return
	}

	// Transfer only up to remaining quota (also clears expires_at)
	err = s.Repo.TransferAnonymousURLsToUserWithLimit(ctx, db.TransferAnonymousURLsToUserWithLimitParams{
		OwnerID:   anonID,
		OwnerID_2: userID,
		Limit:     int32(remaining),
	})
	if err != nil {
		s.Logger.Error("failed to transfer anonymous urls", "anon_id", anonID, "user_id", userID, "error", err)
		return
	}

	s.Logger.Info("transferred anonymous urls with quota limit", "anon_id", anonID, "user_id", userID, "max_transferred", remaining)
}
