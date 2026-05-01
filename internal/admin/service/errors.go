package service

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid admin credentials")
	ErrTokenCreation      = errors.New("failed to create token")
)
