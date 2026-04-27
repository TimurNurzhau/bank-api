package handlers

import (
	"encoding/json"
	"net/http"

	"bank-api/middleware"
	"bank-api/models"
	"bank-api/response"
	"bank-api/services"
	"bank-api/utils"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type CardHandler struct {
	cardService *services.CardService
	validator   *validator.Validate
	logger      *logrus.Logger
}

func NewCardHandler(cardService *services.CardService, logger *logrus.Logger) *CardHandler {
	return &CardHandler{
		cardService: cardService,
		validator:   validator.New(),
		logger:      logger,
	}
}

func (h *CardHandler) Issue(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req struct {
		AccountID int `json:"account_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	card, err := h.cardService.IssueCard(userID, req.AccountID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(card); err != nil {
		h.logger.WithError(err).Error("failed to encode response")
	}
}

func (h *CardHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	cards, err := h.cardService.GetUserCards(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Дополнительная проверка: маскируем номера карт для безопасности
	for i := range cards {
		cards[i].MaskedNumber = utils.MaskCardNumber(cards[i].MaskedNumber)
	}

	if cards == nil {
		cards = []models.Card{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cards); err != nil {
		h.logger.WithError(err).Error("failed to encode response")
	}
}

func (h *CardHandler) Pay(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req struct {
		CardID int     `json:"card_id" validate:"required"`
		Amount float64 `json:"amount" validate:"required,gt=0"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Проверка прав происходит внутри PayWithCard через FindByIDAndUserID
	if err := h.cardService.PayWithCard(req.CardID, userID, req.Amount); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Payment completed successfully",
	}); err != nil {
		h.logger.WithError(err).Error("failed to encode response")
	}
}