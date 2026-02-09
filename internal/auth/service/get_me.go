package service

import (
	"context"

	"url_shortner_backend/pkg/tier"
)

func (s *AuthSvcImp) GetMe(ctx context.Context, userID string) (MeOutput, error) {
	user, err := s.Repo.GetUserByID(ctx, stringToUUID(userID))
	if err != nil {
		s.Logger.Error("failed to get user", "user_id", userID, "error", err)
		return MeOutput{}, ErrUserNotFound
	}

	// Get URL count for this month
	urlsThisMonth, err := s.ShortURLSvc.GetUserURLCountThisMonth(ctx, userID)
	if err != nil {
		s.Logger.Error("failed to get url count", "user_id", userID, "error", err)
		urlsThisMonth = 0
	}

	userTier := tier.Tier(user.Tier)
	limits := tier.GetLimits(userTier)

	output := MeOutput{
		UserID:               uuidToString(user.ID),
		Email:                user.Email,
		Tier:                 string(user.Tier),
		URLsCreatedThisMonth: urlsThisMonth,
		URLsLimit:            limits.URLsPerMonth,
	}

	if user.Name.Valid {
		output.Name = user.Name.String
	}
	if user.AvatarUrl.Valid {
		output.AvatarURL = user.AvatarUrl.String
	}
	if user.SubscriptionEndsAt.Valid {
		t := user.SubscriptionEndsAt.Time
		output.SubscriptionEndsAt = &t
	}

	return output, nil
}
