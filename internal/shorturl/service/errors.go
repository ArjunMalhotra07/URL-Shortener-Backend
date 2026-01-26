package service

import "errors"

var (
	ErrInvalidURL    = errors.New("invalid url provided")
	ErrInvalidCode   = errors.New("invalid code provided")
	ErrInvalidOwner  = errors.New("invalid owner id provided")
	ErrURLCreation   = errors.New("failed to create short url")
	ErrURLCodeUpdate = errors.New("failed to update short url code")
	ErrURLNotFound   = errors.New("short url not found")
	ErrURLExpired    = errors.New("short url has expired")
	ErrURLInactive   = errors.New("short url is inactive")
	ErrURLFetch      = errors.New("failed to fetch urls")
	ErrURLTransfer   = errors.New("failed to transfer urls")
)
