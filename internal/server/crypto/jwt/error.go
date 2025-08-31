// Package jwt предоставляет функционал для работы с JWT токенами.
package jwt

import (
	"errors"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

// Возможные ошибки при работе с токеном.
var (
	ErrTokenInvalid              = errors.New("invalid token")
	ErrTokenInvalidClaims        = jwtlib.ErrTokenInvalidClaims
	ErrTokenInvalidFormat        = errors.New("invalid token format")
	ErrTokenRequiredClaimMissing = jwtlib.ErrTokenRequiredClaimMissing
)
