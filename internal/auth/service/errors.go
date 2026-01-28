package service

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailExists        = errors.New("email already exists")
	ErrUserCreation       = errors.New("failed to create user")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrTokenCreation      = errors.New("failed to create token")
	ErrUserNotFound       = errors.New("user not found")
)
