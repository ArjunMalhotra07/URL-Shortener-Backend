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
		s.Logger.Error().Msg("empty anon_id or user_id")
		return ErrInvalidOwner
	}

	err := s.Repo.TransferAnonymousURLsToUser(ctx, db.TransferAnonymousURLsToUserParams{
		OwnerID:   input.AnonID,
		OwnerID_2: input.UserID,
	})
	if err != nil {
		s.Logger.Err(err).Str("anon_id", input.AnonID).Str("user_id", input.UserID).Msg("failed to transfer urls")
		return ErrURLTransfer
	}

	s.Logger.Info().Str("anon_id", input.AnonID).Str("user_id", input.UserID).Msg("urls transferred")
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
		s.Logger.Err(err).Str("user_id", userID).Msg("failed to count user urls this month")
		return
	}

	remaining := s.Cfg.MonthlyQuotaUser - int(monthCount)
	if remaining <= 0 {
		s.Logger.Info().Str("user_id", userID).Int64("month_count", monthCount).Msg("user at monthly quota, skipping anonymous url transfer")
		return
	}

	// Transfer only up to remaining quota (also clears expires_at)
	err = s.Repo.TransferAnonymousURLsToUserWithLimit(ctx, db.TransferAnonymousURLsToUserWithLimitParams{
		OwnerID:   anonID,
		OwnerID_2: userID,
		Limit:     int32(remaining),
	})
	if err != nil {
		s.Logger.Err(err).Str("anon_id", anonID).Str("user_id", userID).Msg("failed to transfer anonymous urls")
		return
	}

	s.Logger.Info().Str("anon_id", anonID).Str("user_id", userID).Int("max_transferred", remaining).Msg("transferred anonymous urls with quota limit")
}
