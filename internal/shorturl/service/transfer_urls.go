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
