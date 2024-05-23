package domain

import "errors"

var (
	ErrTokenIsNotValid         = errors.New("token is not valid")
	ErrInvalidTokenClaims      = errors.New("invalid token claims")
	ErrUserIDClaimNotFound     = errors.New("user_id claim not found or invalid")
	ErrTokenIsExpired          = errors.New("token is expired")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
)
