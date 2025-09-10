package jwt_test

import (
	"testing"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/mr-filatik/go-goph-keeper/internal/server/crypto/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
	===== NewEncryptor =====
*/

func TestNewEncryptor(t *testing.T) {
	t.Parallel()

	encryptorOpts := jwt.WithExpireTime(15 * time.Minute)

	encryptor := jwt.NewEncryptor("TEST_SECRET_KEY", encryptorOpts)

	assert.NotEmpty(t, encryptor)
}

/*
	===== Encryptor.CreateClaimsWithUserID =====
*/

func TestEncryptor_CreateClaimsWithUserID(t *testing.T) {
	t.Parallel()

	encryptorOpts := jwt.WithExpireTime(15 * time.Minute)

	encryptor := jwt.NewEncryptor("TEST_SECRET_KEY", encryptorOpts)

	claims := encryptor.CreateClaimsWithUserID("test_user_id")

	assert.Equal(t, "test_user_id", claims["user_id"])
}

/*
	===== Encryptor.GenerateTokenString =====
*/

func TestEncryptor_GenerateTokenString(t *testing.T) {
	t.Parallel()

	encryptorOpts := jwt.WithExpireTime(1 * time.Second)

	encryptor := jwt.NewEncryptor("TEST_SECRET_KEY", encryptorOpts)

	claims := encryptor.CreateClaimsWithUserID("test_user_id")

	tokenString, generErr := encryptor.GenerateTokenString(claims)
	require.NoError(t, generErr)

	token, parseErr := jwtlib.Parse(tokenString, func(_ *jwtlib.Token) (interface{}, error) {
		return []byte("TEST_SECRET_KEY"), nil
	})
	require.NoError(t, parseErr)

	assert.True(t, token.Valid)

	time.Sleep(1 * time.Second)

	token, parseErr = jwtlib.Parse(tokenString, func(_ *jwtlib.Token) (interface{}, error) {
		return []byte("TEST_SECRET_KEY"), nil
	})
	require.ErrorIs(t, parseErr, jwtlib.ErrTokenExpired)

	assert.False(t, token.Valid)
}

/*
	===== Encryptor.ValidateTokenBearer =====
*/

func TestEncryptor_ValidateTokenBearer(t *testing.T) {
	t.Parallel()

	encryptor := jwt.NewEncryptor("TEST_SECRET_KEY")

	claims := encryptor.CreateClaimsWithUserID("test_user_id")

	tokenString, generErr := encryptor.GenerateTokenString(claims)
	require.NoError(t, generErr)

	token, parseErr := encryptor.ValidateTokenBearer(tokenString)
	require.ErrorIs(t, parseErr, jwt.ErrTokenInvalidFormat)
	require.Nil(t, token)

	token, parseErr = encryptor.ValidateTokenBearer("Bearer " + tokenString)
	require.NoError(t, parseErr)
	require.NotNil(t, token)
	assert.True(t, token.Valid)
}

/*
	===== Encryptor.ValidateToken =====
*/

func TestEncryptor_ValidateToken(t *testing.T) {
	t.Parallel()

	encryptorOpts := jwt.WithExpireTime(1 * time.Second)

	encryptor := jwt.NewEncryptor("TEST_SECRET_KEY", encryptorOpts)

	claims := encryptor.CreateClaimsWithUserID("test_user_id")

	tokenString, generErr := encryptor.GenerateTokenString(claims)
	require.NoError(t, generErr)

	token, parseErr := encryptor.ValidateToken(tokenString)
	require.NoError(t, parseErr)
	require.NotNil(t, token)

	assert.True(t, token.Valid)

	time.Sleep(1 * time.Second)

	token, parseErr = encryptor.ValidateToken(tokenString)
	require.ErrorIs(t, parseErr, jwtlib.ErrTokenExpired)
	require.Nil(t, token)
}

/*
	===== Encryptor.GetClaimUserIDFromToken =====
*/

func TestEncryptor_GetClaimUserIDFromToken(t *testing.T) {
	t.Parallel()

	encryptor := jwt.NewEncryptor("TEST_SECRET_KEY")

	claims := encryptor.CreateClaimsWithUserID("test_user_id")

	tokenString, generErr := encryptor.GenerateTokenString(claims)
	require.NoError(t, generErr)

	token, parseErr := encryptor.ValidateToken(tokenString)
	require.NoError(t, parseErr)
	require.NotNil(t, token)

	user, getErr := encryptor.GetClaimUserIDFromToken(token)
	require.NoError(t, getErr)
	assert.Equal(t, "test_user_id", user)

	claims = jwtlib.MapClaims{}

	tokenString, generErr = encryptor.GenerateTokenString(claims)
	require.NoError(t, generErr)

	token, parseErr = encryptor.ValidateToken(tokenString)
	require.NoError(t, parseErr)
	require.NotNil(t, token)

	_, getErr = encryptor.GetClaimUserIDFromToken(token)
	require.ErrorIs(t, getErr, jwt.ErrTokenRequiredClaimMissing)
}
