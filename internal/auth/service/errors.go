package service

import "errors"

var (
	ErrInvalidCredentials      = errors.New("invalid email or password")
	ErrEmailExists             = errors.New("email already exists")
	ErrEmailExistsWithPassword = errors.New("email already registered with password login")
	ErrEmailExistsWithGoogle   = errors.New("email already registered with google login")
	ErrUserCreation            = errors.New("failed to create user")
	ErrInvalidToken            = errors.New("invalid token")
	ErrTokenExpired            = errors.New("token expired")
	ErrTokenCreation           = errors.New("failed to create token")
	ErrUserNotFound            = errors.New("user not found")
)
