// Package jwt предоставляет функционал для работы с JWT токенами.
package jwt

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// IEncryptor - интерфейс для работы с токенами.
type IEncryptor interface {
	CreateClaimsWithUserID(userID string) jwt.MapClaims
	GenerateTokenString(claims jwt.MapClaims) (string, error)
	ValidateTokenBearer(tokenString string) (*jwt.Token, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	GetClaimUserIDFromToken(token *jwt.Token) (string, error)
}

// Encryptor содержит общие данные для всех хендлеров.
type Encryptor struct {
	tockenExpireTime time.Duration
	secretJWTKey     []byte
}

// EncryptorOption представляет дополнительные опции для Encryptor.
type EncryptorOption func(*Encryptor)

// _defaultTokenExpireTime значение по умолчанию для времени жизни токена.
const defaultTokenExpireTime = 24 * time.Hour

// WithExpireTime устанавливает время жизни токена.
func WithExpireTime(exp time.Duration) EncryptorOption {
	return func(e *Encryptor) {
		e.tockenExpireTime = exp
	}
}

// NewEncryptor создаёт и инициализирует новый экзепляр *Encryptor.
//
// Параметры:
//   - jwtKey: логгер;
//   - exp: логгер.
func NewEncryptor(jwtKey string, opts ...EncryptorOption) *Encryptor {
	encryptor := &Encryptor{
		secretJWTKey:     []byte(jwtKey),
		tockenExpireTime: defaultTokenExpireTime,
	}

	for index := range opts {
		opts[index](encryptor)
	}

	return encryptor
}

// CreateClaimsWithUserID создаёт и инициализирует Claims с user_id.
//
// Параметры:
//   - userId: идентификатор пользователя.
func (e *Encryptor) CreateClaimsWithUserID(userID string) jwt.MapClaims {
	return jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(e.tockenExpireTime).Unix(),
	}
}

// GenerateTokenString создаёт подписанный токен в виде строки.
//
// Параметры:
//   - claims: набор Claims.
func (e *Encryptor) GenerateTokenString(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(e.secretJWTKey)
	if err != nil {
		return "", fmt.Errorf("generate toker: %w", err)
	}

	return tokenStr, nil
}

// ValidateTokenBearer проверяет токен в формате Bearer и выводит его Claims.
//
// Параметры:
//   - tokenString: строка с Bearer токеном.
func (e *Encryptor) ValidateTokenBearer(tokenString string) (*jwt.Token, error) {
	if !strings.HasPrefix(strings.ToLower(tokenString), "bearer ") {
		return nil, fmt.Errorf("no bearer format: %w", ErrTokenInvalidFormat)
	}

	tokenStr := strings.TrimSpace(tokenString[len("Bearer "):])

	return e.ValidateToken(tokenStr)
}

// ValidateToken проверяет токен и выводит его Claims.
//
// Параметры:
//   - tokenString: строка с токеном.
func (e *Encryptor) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, parseErr := jwt.Parse(tokenString, func(_ *jwt.Token) (interface{}, error) {
		return e.secretJWTKey, nil
	})

	if parseErr != nil {
		return nil, fmt.Errorf("parse token: %w", parseErr)
	}

	return token, nil
}

// GetClaimUserIDFromToken получает user_id из Claims токена.
//
// Параметры:
//   - token: токен.
func (e *Encryptor) GetClaimUserIDFromToken(token *jwt.Token) (string, error) {
	claims, claimOk := token.Claims.(jwt.MapClaims)
	if !claimOk {
		return "", fmt.Errorf("map claims: %w", ErrTokenInvalidClaims)
	}

	userID, userOk := claims["user_id"].(string)
	if !userOk {
		return "", fmt.Errorf("user_id: %w", ErrTokenRequiredClaimMissing)
	}

	return userID, nil
}
