// Package auth предоставляет функционал для обработчиков запросов для авторизации.
package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mr-filatik/go-goph-keeper/internal/server/crypto/jwt"
	"github.com/mr-filatik/go-goph-keeper/internal/server/handler"
	"github.com/mr-filatik/go-goph-keeper/internal/server/storage"
	"github.com/mr-filatik/go-goph-keeper/internal/server/storage/entity"
	"golang.org/x/crypto/bcrypt"
)

// Handler хранит данные необходимые для обработчиков.
type Handler struct {
	handler.Handler
	encryptor *jwt.Encryptor
}

// HandlerOption представляет дополнительные опции для Handler.
type HandlerOption func(*Handler)

// func WithLogger(l *slog.Logger) Option { return func(h *Handler){ h.log = l } }

// NewHandler создаёт новый экземпляр Handler.
func NewHandler(hand handler.Handler, enc *jwt.Encryptor, opts ...HandlerOption) *Handler {
	authHandler := &Handler{
		Handler:   hand,
		encryptor: enc,
	}

	for index := range opts {
		opts[index](authHandler)
	}

	return authHandler
}

// UserRegister регистрирует нового пользователя.
func (h *Handler) UserRegister(writer http.ResponseWriter, req *http.Request) {
	var data registerReq

	err := handler.GetDataFromBodyJSON(req, &data)
	if err != nil {
		h.ResponseError(writer, http.StatusBadRequest, err)

		return
	}

	passHash, hashErr := generatePasswordHash(data.Password)
	if hashErr != nil {
		h.ResponseError(writer, http.StatusInternalServerError, hashErr)

		return
	}

	user := entity.NewUser(data.Email, passHash)

	_, addErr := h.Stor.AddNewUser(req.Context(), user)
	if addErr != nil {
		if errors.Is(addErr, storage.ErrEntityAlreadyExists) {
			h.ResponseError(writer, http.StatusConflict, addErr)

			return
		}

		h.ResponseError(writer, http.StatusInternalServerError, addErr)

		return
	}

	claims := h.encryptor.CreateClaimsWithUserID(user.ID)

	token, tokenErr := h.encryptor.GenerateTokenString(claims)
	if tokenErr != nil {
		h.ResponseError(writer, http.StatusInternalServerError, tokenErr)

		return
	}

	_, atErr := h.Stor.AddNewToken(req.Context(), user.ID, &entity.Token{})
	if atErr != nil {
		h.ResponseError(writer, http.StatusInternalServerError, atErr)

		return
	}

	h.ResponceWithJSON(writer, registerResp{
		Token: token,
	})
}

// UserLogin авторизует нового пользователя.
func (h *Handler) UserLogin(writer http.ResponseWriter, req *http.Request) {
	var data loginReq

	dataErr := handler.GetDataFromBodyJSON(req, &data)
	if dataErr != nil {
		h.ResponseError(writer, http.StatusBadRequest, dataErr)

		return
	}

	user, findErr := h.Stor.FindUserByEmail(req.Context(), data.Email)
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

	claims := h.encryptor.CreateClaimsWithUserID(user.ID)

	token, tokenErr := h.encryptor.GenerateTokenString(claims)
	if tokenErr != nil {
		h.ResponseError(writer, http.StatusInternalServerError, tokenErr)

		return
	}

	_, atErr := h.Stor.AddNewToken(req.Context(), user.ID, &entity.Token{})
	if atErr != nil {
		h.ResponseError(writer, http.StatusInternalServerError, atErr)

		return
	}

	h.ResponceWithJSON(writer, loginResp{
		Token: token,
	})
}

// UserLogout убирает авторизацию для пользователя.
func (h *Handler) UserLogout(writer http.ResponseWriter, req *http.Request) {
	authToken := req.Header.Get("Authorization")

	token, err := h.encryptor.ValidateTokenBearer(authToken)
	if err != nil {
		h.ResponseError(writer, http.StatusUnauthorized, err)

		return
	}

	userID, userErr := h.encryptor.GetClaimUserIDFromToken(token)
	if userErr != nil {
		h.ResponseError(writer, http.StatusBadRequest, userErr)

		return
	}

	_ = userID
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

	return err == nil
}
