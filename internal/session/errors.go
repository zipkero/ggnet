package session

import "errors"

var (
	ErrSessionClosed       = errors.New("session closed")
	ErrSessionNotStarted   = errors.New("session not started")
	ErrSessionNotConnected = errors.New("session not connected")
	ErrSessionNotStopped   = errors.New("session not stopped")
	ErrSessionNotFound     = errors.New("session not found")
)
