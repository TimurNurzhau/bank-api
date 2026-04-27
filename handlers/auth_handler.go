package handlers

import (
	"encoding/json"
	"net/http"

	"bank-api/models"
	"bank-api/repositories"
	"bank-api/response"
	"bank-api/services"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	authService *services.AuthService
	validator   *validator.Validate
	logger      *logrus.Logger
}

func NewAuthHandler(authService *services.AuthService, logger *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator.New(),
		logger:      logger,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.authService.Register(&req)
	if err != nil {
		// Обработка понятных ошибок
		switch err {
		case repositories.ErrDuplicateUsername:
			response.Error(w, http.StatusConflict, "username already taken")
			return
		case repositories.ErrDuplicateEmail:
			response.Error(w, http.StatusConflict, "email already registered")
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.WithError(err).Error("failed to encode response")
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.WithError(err).Error("failed to encode response")
	}
}