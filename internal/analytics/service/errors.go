package service

import "errors"

var (
	ErrURLNotFound     = errors.New("short url not found")
	ErrURLNotOwned     = errors.New("url not owned by user")
	ErrAnalyticsFetch  = errors.New("failed to fetch analytics")
	ErrClickRecording  = errors.New("failed to record click")
	ErrInvalidTimeRange = errors.New("invalid time range")
)
