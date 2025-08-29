package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mr-filatik/go-goph-keeper/internal/server/handler"
	"github.com/mr-filatik/go-goph-keeper/internal/server/storage"
	"github.com/mr-filatik/go-goph-keeper/internal/server/storage/entity"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler хранит данные необходимые для обработчиков.
type AuthHandler struct {
	handler.Handler
}

// AuthHandlerOption представляет дополнительные опции для Handler.
type AuthHandlerOption func(*AuthHandler)

// func WithLogger(l *slog.Logger) Option { return func(h *Handler){ h.log = l } }

// NewHandler создаёт новый экземпляр Handler.
func NewHandler(hand handler.Handler, opts ...AuthHandlerOption) *AuthHandler {
	h := &AuthHandler{
		Handler: hand,
	}

	for _, o := range opts {
		o(h)
	}

	return h
}

// UserRegister регистрирует нового пользователя.
func (h *AuthHandler) UserRegister(writer http.ResponseWriter, req *http.Request) {
	var data registerReq
	err := handler.GetDataFromBodyJSON(req, &data)
	if err != nil {
		h.ResponseError(writer, http.StatusBadRequest, err)

		return
	}

	passHash, hashErr := generatePasswordHash(data.Password)
	if hashErr != nil {
		h.ResponseError(writer, http.StatusInternalServerError, err)

		return
	}

	user := entity.User{
		Email:        data.Email,
		PasswordHash: passHash,
	}

	_, addErr := h.Stor.AddNewUser(context.Background(), &user)
	if addErr != nil {
		if errors.Is(addErr, storage.ErrEntityAlreadyExists) {
			h.ResponseError(writer, http.StatusConflict, addErr)

			return
		}

		h.ResponseError(writer, http.StatusInternalServerError, addErr)

		return
	}

	token, tokenErr := generateToken(user.ID)
	if tokenErr != nil {
		h.ResponseError(writer, http.StatusInternalServerError, err)

		return
	}

	_, atErr := h.Stor.AddNewToken(context.Background(), user.ID, &entity.Token{})
	if atErr != nil {
		h.ResponseError(writer, http.StatusInternalServerError, err)

		return
	}

	resp := registerResp{
		Token: token,
	}

	h.ResponceWithJSON(writer, resp)
}

// UserLogin авторизует нового пользователя.
func (h *AuthHandler) UserLogin(writer http.ResponseWriter, req *http.Request) {
	var data loginReq
	err := handler.GetDataFromBodyJSON(req, &data)
	if err != nil {
		h.ResponseError(writer, http.StatusBadRequest, err)

		return
	}

	user, findErr := h.Stor.FindUserByEmail(context.Background(), data.Email)
	if findErr != nil {
		if errors.Is(findErr, storage.ErrEntityNotFound) {
			h.ResponseError(writer, http.StatusNotFound, findErr)

			return
		}

		h.ResponseError(writer, http.StatusInternalServerError, findErr)

		return
	}

	if ok := comparePasswordHash(data.Password, user.PasswordHash); !ok {
		h.ResponseError(writer, http.StatusUnauthorized, findErr)

		return
	}

	token, tokenErr := generateToken(user.ID)
	if tokenErr != nil {
		h.ResponseError(writer, http.StatusInternalServerError, err)

		return
	}

	_, atErr := h.Stor.AddNewToken(context.Background(), user.ID, &entity.Token{})
	if atErr != nil {
		h.ResponseError(writer, http.StatusInternalServerError, err)

		return
	}

	resp := registerResp{
		Token: token,
	}

	h.ResponceWithJSON(writer, resp)
}

// UserLogout убирает авторизацию для пользователя.
func (h *AuthHandler) UserLogout(writer http.ResponseWriter, req *http.Request) {
	head := req.Header.Get("Authorization")

	if !strings.HasPrefix(strings.ToLower(head), "bearer ") {
		h.ResponseError(writer, http.StatusUnauthorized, fmt.Errorf("missing bearer token"))

		return
	}

	tokenStr := strings.TrimSpace(head[len("Bearer "):])

	token, err := parseToken(tokenStr)
	if err != nil {
		h.ResponseError(writer, http.StatusUnauthorized, err)
	}

	if !token.Valid {
		h.ResponseError(writer, http.StatusUnauthorized, fmt.Errorf("bearer token not valid"))
	}

	//token.Claims.
}

func generatePasswordHash(password string) (string, error) {
	bytePass := []byte(password)

	hash, err := bcrypt.GenerateFromPassword(bytePass, bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("generate password: %w", err)
	}

	return string(hash), nil
}

func comparePasswordHash(password, hash string) bool {
	bytePass := []byte(password)
	byteHash := []byte(hash)

	err := bcrypt.CompareHashAndPassword(byteHash, bytePass)

	if err != nil {
		return false
	}

	return true
}

var secretKey = []byte("my_secret_key")

func generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("generate toker: %w", err)
	}

	return tokenStr, nil
}

func parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
}
