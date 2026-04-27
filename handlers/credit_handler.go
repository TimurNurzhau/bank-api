package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bank-api/middleware"
	"bank-api/models"
	"bank-api/response"
	"bank-api/services"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type CreditHandler struct {
	creditService *services.CreditService
	cbrService    *services.CBRService
	logger        *logrus.Logger
}

func NewCreditHandler(creditService *services.CreditService, cbrService *services.CBRService, logger *logrus.Logger) *CreditHandler {
	return &CreditHandler{
		creditService: creditService,
		cbrService:    cbrService,
		logger:        logger,
	}
}

func (h *CreditHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req models.CreateCreditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	keyRate, err := h.cbrService.GetKeyRate()
	if err != nil {
		response.Error(w, http.StatusServiceUnavailable, "failed to get key rate from CBR: "+err.Error())
		return
	}

	credit, err := h.creditService.CreateCredit(userID, &req, keyRate)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(credit); err != nil {
		h.logger.WithError(err).Error("failed to encode response")
	}
}

func (h *CreditHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	credits, err := h.creditService.GetUserCredits(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(credits); err != nil {
		h.logger.WithError(err).Error("failed to encode response")
	}
}

func (h *CreditHandler) GetSchedule(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	vars := mux.Vars(r)
	creditID, err := strconv.Atoi(vars["creditId"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid credit id")
		return
	}

	schedule, err := h.creditService.GetCreditSchedule(creditID, userID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(schedule); err != nil {
		h.logger.WithError(err).Error("failed to encode response")
	}
}

func (h *CreditHandler) EarlyRepayment(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	vars := mux.Vars(r)
	creditID, err := strconv.Atoi(vars["creditId"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid credit id")
		return
	}

	var req struct {
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Проверка amount > 0 выполняется внутри creditService.EarlyRepayment
	// Не дублируем проверку здесь

	if err := h.creditService.EarlyRepayment(creditID, userID, req.Amount); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Early repayment completed",
	}); err != nil {
		h.logger.WithError(err).Error("failed to encode response")
	}
}