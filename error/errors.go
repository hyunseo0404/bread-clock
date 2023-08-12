package errors

import "errors"

var (
	ErrAuthInvalidProvider = errors.New("invalid provider")
	ErrAuthInvalidToken    = errors.New("invalid token")
)
