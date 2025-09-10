package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mr-filatik/go-goph-keeper/internal/server/crypto/jwt"
	"github.com/mr-filatik/go-goph-keeper/internal/server/handler"
	"github.com/mr-filatik/go-goph-keeper/internal/server/handler/auth"
	"github.com/mr-filatik/go-goph-keeper/internal/server/storage"
	"github.com/mr-filatik/go-goph-keeper/internal/server/storage/entity"
	"github.com/mr-filatik/go-goph-keeper/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

/*
	===== NewHandler =====
*/

func TestNewHandler(t *testing.T) {
	t.Parallel()

	mockLogger := testutil.NewMockLogger()

	mainHandler := handler.NewHandler(nil, mockLogger)
	authHandler := auth.NewHandler(*mainHandler, nil)

	assert.NotEmpty(t, authHandler)
}

/*
	===== Handler.UserRegister =====
*/

type argUserRegister struct {
	body map[string]string
}

type wantUserRegister struct {
	body       string
	statusCode int
}

type testUserRegister struct {
	name string
	args argUserRegister
	want wantUserRegister
}

func createTestsForUserRegister() []testUserRegister {
	tests := []testUserRegister{
		{
			name: "correct user",
			args: argUserRegister{
				body: map[string]string{
					"email":    "test@example.com",
					"password": "P@ssw0rd!",
				},
			},
			want: wantUserRegister{
				statusCode: http.StatusOK,
				body:       "",
			},
		},
		{
			name: "uncorrect user",
			args: argUserRegister{
				body: map[string]string{
					"email":    "register-user",
					"password": "P@ssw0rd!",
				},
			},
			want: wantUserRegister{
				statusCode: http.StatusConflict,
				body:       "Error\n",
			},
		},
	}

	return tests
}

func TestHandler_UserRegister(t *testing.T) {
	t.Parallel()

	mockLogger := testutil.NewMockLogger()

	mockEncryptor := jwt.NewEncryptor("TEST_SECRET_KEY")

	mockStore := &mockStorage{
		addNewUserFn: func(_ context.Context, user *entity.User) (string, error) {
			if user.Email == "register-user" {
				return "", storage.ErrEntityAlreadyExists
			}

			return user.ID, nil
		},
		addNewTokenFn: func(_ context.Context, userID string, _ *entity.Token) (string, error) {
			if userID == "register-user-id" {
				return "", storage.ErrEntityAlreadyExists
			}

			return userID, nil
		},
	}

	mainHandler := handler.NewHandler(mockStore, mockLogger)
	authHandler := auth.NewHandler(*mainHandler, mockEncryptor)

	require.NotEmpty(t, authHandler)

	tests := createTestsForUserRegister()

	for index := range tests {
		internalTest := tests[index]
		t.Run(internalTest.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			if internalTest.args.body != nil {
				_ = json.NewEncoder(&buf).Encode(internalTest.args.body)
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/register", &buf)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			authHandler.UserRegister(recorder, req)

			AnalizeResponse(t,
				internalTest.want.statusCode, recorder.Code,
				internalTest.want.body, recorder.Body.String(),
			)
		})
	}
}

/*
	===== Handler.UserLogin =====
*/

type argUserLogin struct {
	body map[string]string
}

type wantUserLogin struct {
	body       string
	statusCode int
}

type testUserLogin struct {
	name string
	args argUserLogin
	want wantUserLogin
}

func createTestsForUserLogin() []testUserLogin {
	tests := []testUserLogin{
		{
			name: "correct user",
			args: argUserLogin{
				body: map[string]string{
					"email":    "test@example.com",
					"password": "password-hash",
				},
			},
			want: wantUserLogin{
				statusCode: http.StatusOK,
				body:       "",
			},
		},
		{
			name: "correct user with uncorrect password",
			args: argUserLogin{
				body: map[string]string{
					"email":    "test@example.com",
					"password": "password",
				},
			},
			want: wantUserLogin{
				statusCode: http.StatusUnauthorized,
				body:       "Error\n",
			},
		},
		{
			name: "uncorrect user",
			args: argUserLogin{
				body: map[string]string{
					"email":    "user",
					"password": "password",
				},
			},
			want: wantUserLogin{
				statusCode: http.StatusNotFound,
				body:       "Error\n",
			},
		},
	}

	return tests
}

func TestHandler_UserLogin(t *testing.T) {
	t.Parallel()

	mockLogger := testutil.NewMockLogger()
	mockEncryptor := jwt.NewEncryptor("TEST_SECRET_KEY")

	bytePass := []byte("password-hash")
	hash, _ := bcrypt.GenerateFromPassword(bytePass, bcrypt.DefaultCost)

	mockStore := &mockStorage{
		findUserByEmailFn: func(_ context.Context, email string) (*entity.User, error) {
			if email == "user" {
				return nil, storage.ErrEntityNotFound
			}

			return &entity.User{
				ID:           "login-user-id",
				Email:        email,
				PasswordHash: string(hash),
			}, nil
		},
		addNewTokenFn: func(_ context.Context, userID string, _ *entity.Token) (string, error) {
			if userID != "login-user-id" {
				return "", storage.ErrEntityAlreadyExists
			}

			return userID, nil
		},
	}

	mainHandler := handler.NewHandler(mockStore, mockLogger)
	authHandler := auth.NewHandler(*mainHandler, mockEncryptor)

	require.NotEmpty(t, authHandler)

	tests := createTestsForUserLogin()

	for index := range tests {
		internalTest := tests[index]
		t.Run(internalTest.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			if internalTest.args.body != nil {
				_ = json.NewEncoder(&buf).Encode(internalTest.args.body)
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/login", &buf)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			authHandler.UserLogin(recorder, req)

			AnalizeResponse(t,
				internalTest.want.statusCode, recorder.Code,
				internalTest.want.body, recorder.Body.String(),
			)
		})
	}
}

/*
	===== Handler.UserLogout =====
*/

type argUserLogout struct {
	userID string
}

type wantUserLogout struct {
	body       string
	statusCode int
}

type testUserLogout struct {
	name string
	args argUserLogout
	want wantUserLogout
}

func createTestsForUserLogout() []testUserLogout {
	tests := []testUserLogout{
		{
			name: "correct user",
			args: argUserLogout{
				userID: "user-id",
			},
			want: wantUserLogout{
				statusCode: http.StatusOK,
				body:       "",
			},
		},
		{
			name: "uncorrect user",
			args: argUserLogout{
				userID: "uncorrect-user-id",
			},
			want: wantUserLogout{
				statusCode: http.StatusUnauthorized,
				body:       "Error\n",
			},
		},
	}

	return tests
}

func TestHandler_UserLogout(t *testing.T) {
	t.Parallel()

	mockLogger := testutil.NewMockLogger()

	mockEncryptor := jwt.NewEncryptor("TEST_SECRET_KEY")

	mockStore := &mockStorage{
		isTokenByUserIDFn: func(_ context.Context, userID string) bool {
			return userID == "user-id"
		},
		deleteTokenFn: func(_ context.Context, _ string) error {
			return nil
		},
	}

	mainHandler := handler.NewHandler(mockStore, mockLogger)
	authHandler := auth.NewHandler(*mainHandler, mockEncryptor)

	require.NotEmpty(t, authHandler)

	tests := createTestsForUserLogout()

	for index := range tests {
		internalTest := tests[index]
		t.Run(internalTest.name, func(t *testing.T) {
			t.Parallel()

			// Сделать генерацию токена для Bearer ... в заголовке Authorization
			claims := mockEncryptor.CreateClaimsWithUserID(internalTest.args.userID)
			token, _ := mockEncryptor.GenerateTokenString(claims)

			req := httptest.NewRequest(http.MethodPost, "/auth/login", http.NoBody)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", strings.Join([]string{"Bearer", token}, " "))

			recorder := httptest.NewRecorder()

			authHandler.UserLogout(recorder, req)

			assert.Equal(t, internalTest.want.statusCode, recorder.Code)

			assert.Equal(t, internalTest.want.body, recorder.Body.String())
		})
	}
}

/*
	===== Helpers =====
*/

func AnalizeResponse(t *testing.T, expStat, actStat int, expBody, actBody string) {
	t.Helper()

	assert.Equal(t, expStat, actStat)

	if actStat == http.StatusOK {
		assert.NotEmpty(t, actBody)
	} else {
		assert.Equal(t, expBody, actBody)
	}
}

/*
	===== Mock IStorage =====
*/

type mockStorage struct {
	addNewUserFn      func(ctx context.Context, user *entity.User) (string, error)
	findUserByEmailFn func(ctx context.Context, email string) (*entity.User, error)
	addNewTokenFn     func(ctx context.Context, userID string, token *entity.Token) (string, error)
	isTokenByUserIDFn func(ctx context.Context, userID string) bool
	deleteTokenFn     func(ctx context.Context, userID string) error
}

func (m *mockStorage) AddNewUser(ctx context.Context, user *entity.User) (string, error) {
	return m.addNewUserFn(ctx, user)
}

func (m *mockStorage) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return m.findUserByEmailFn(ctx, email)
}

func (m *mockStorage) AddNewToken(
	ctx context.Context,
	userID string,
	token *entity.Token,
) (string, error) {
	return m.addNewTokenFn(ctx, userID, token)
}

func (m *mockStorage) IsTokenByUserID(ctx context.Context, userID string) bool {
	return m.isTokenByUserIDFn(ctx, userID)
}

func (m *mockStorage) DeleteToken(ctx context.Context, userID string) error {
	return m.deleteTokenFn(ctx, userID)
}
