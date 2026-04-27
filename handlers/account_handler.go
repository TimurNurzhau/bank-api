package handlers

import (
	"encoding/json"
	"net/http"

	"bank-api/middleware"
	"bank-api/models"
	"bank-api/response"
	"bank-api/services"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type AccountHandler struct {
	accountService *services.AccountService
	validator      *validator.Validate
	logger         *logrus.Logger
}

func NewAccountHandler(accountService *services.AccountService, logger *logrus.Logger) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
		validator:      validator.New(),
		logger:         logger,
	}
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req models.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	account, err := h.accountService.CreateAccount(userID, &req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(account); err != nil {
		h.logger.WithError(err).Error("failed to encode response")
	}
}

func (h *AccountHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	accounts, err := h.accountService.GetUserAccounts(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Проверка: убеждаемся что все счета принадлежат пользователю
	for _, account := range accounts {
		if account.UserID != userID {
			response.Error(w, http.StatusForbidden, "access denied to some accounts")
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(accounts); err != nil {
		h.logger.WithError(err).Error("failed to encode response")
	}
}