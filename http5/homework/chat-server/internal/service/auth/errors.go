package auth

import "errors"

var (
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrIncorrectTypeClaims  = errors.New("token claims are not of type *tokenClaims")
)
