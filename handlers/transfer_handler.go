package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"bank-api/middleware"
	"bank-api/models"
	"bank-api/response"
	"bank-api/services"

	"github.com/go-playground/validator/v10"
)

type TransferHandler struct {
	transferService *services.TransferService
	validator       *validator.Validate
}

func NewTransferHandler(transferService *services.TransferService) *TransferHandler {
	return &TransferHandler{
		transferService: transferService,
		validator:       validator.New(),
	}
}

func (h *TransferHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req models.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Дополнительная проверка: нельзя переводить на тот же счет
	if req.FromAccountID == req.ToAccountID {
		response.Error(w, http.StatusBadRequest, "cannot transfer to the same account")
		return
	}

	// Проверка прав на from_account происходит внутри TransferService
	if err := h.transferService.Transfer(userID, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "success"}); err != nil {
		log.Printf("encode error: %v", err)
	}
}

func (h *TransferHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req models.DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Проверка прав на account происходит внутри Deposit
	if err := h.transferService.Deposit(userID, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "success"}); err != nil {
		log.Printf("encode error: %v", err)
	}
}