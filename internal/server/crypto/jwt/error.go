// Package jwt предоставляет функционал для работы с JWT токенами.
package jwt

import (
	"errors"

	externaljwt "github.com/golang-jwt/jwt/v5"
)

// Возможные ошибки при работе с токеном.
var (
	ErrTokenInvalid              = errors.New("invalid token")
	ErrTokenInvalidClaims        = externaljwt.ErrTokenInvalidClaims
	ErrTokenInvalidFormat        = errors.New("invalid token format")
	ErrTokenRequiredClaimMissing = externaljwt.ErrTokenRequiredClaimMissing
)
