package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"bank-api/middleware"
	"bank-api/models"
	"bank-api/services"

	"github.com/go-playground/validator/v10"
)

type CardHandler struct {
	cardService *services.CardService
	validator   *validator.Validate
}

func NewCardHandler(cardService *services.CardService) *CardHandler {
	return &CardHandler{
		cardService: cardService,
		validator:   validator.New(),
	}
}

func (h *CardHandler) Issue(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req struct {
		AccountID int `json:"account_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	card, err := h.cardService.IssueCard(userID, req.AccountID)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(card); err != nil {
		log.Printf("encode error: %v", err)
	}
}

func (h *CardHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	cards, err := h.cardService.GetUserCards(userID)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	if cards == nil {
		cards = []models.Card{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cards); err != nil {
		log.Printf("encode error: %v", err)
	}
}

// Оплата картой
func (h *CardHandler) Pay(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req struct {
		CardID int     `json:"card_id" validate:"required"`
		Amount float64 `json:"amount" validate:"required,gt=0"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	if err := h.cardService.PayWithCard(req.CardID, userID, req.Amount); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Payment completed successfully",
	})
}