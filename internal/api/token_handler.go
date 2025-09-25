package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/martialanouman/femProject/internal/store"
	"github.com/martialanouman/femProject/internal/tokens"
	"github.com/martialanouman/femProject/internal/utils"
)

type createTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenHandler struct {
	store     store.TokenStore
	userStore store.UserStore
	logger    *log.Logger
}

func NewTokenHandler(store store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		store:     store,
		userStore: userStore,
		logger:    logger,
	}
}

func (h *TokenHandler) validateCreateTokenRequest(payload createTokenRequest) error {
	if payload.Password == "" {
		return errors.New("password is required")
	}

	if payload.Username == "" {
		return errors.New("username is required")
	}

	return nil
}

func (h *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: decoding request %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	err = h.validateCreateTokenRequest(req)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	user, err := h.userStore.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		h.logger.Printf("ERROR: fetching user %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	ok, err := user.PasswordHash.Matches(req.Password)
	if err != nil {
		h.logger.Printf("ERROR: password matches %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "unauthorized"})
		return
	}

	const expiry24Hour = 24 * time.Hour
	token, err := h.store.CreateToken(user.Id, expiry24Hour, tokens.ScopeAuth)
	if err != nil {
		h.logger.Printf("ERROR: creating token %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"auth_token": token})
}

func (h *TokenHandler) HandleRevokeAllTokensForUser(w http.ResponseWriter, r *http.Request) {

}
