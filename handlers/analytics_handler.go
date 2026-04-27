package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bank-api/middleware"
	"bank-api/response"
	"bank-api/services"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
	logger           *logrus.Logger
}

func NewAnalyticsHandler(analyticsService *services.AnalyticsService, logger *logrus.Logger) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
		logger:           logger,
	}
}

func (h *AnalyticsHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	stats, err := h.analyticsService.GetMonthlyStats(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	creditLoad, _ := h.analyticsService.GetCreditLoad(userID)

	result := map[string]interface{}{
		"monthly_income":   stats["income"],
		"monthly_expenses": stats["expenses"],
		"credit_load":      creditLoad,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		h.logger.WithError(err).Error("failed to encode response")
	}
}

func (h *AnalyticsHandler) PredictBalance(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	vars := mux.Vars(r)

	accountID, err := strconv.Atoi(vars["accountId"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid account id")
		return
	}

	daysStr := r.URL.Query().Get("days")
	days := 30
	if daysStr != "" {
		parsed, err := strconv.Atoi(daysStr)
		if err == nil && parsed > 0 {
			days = parsed
		}
	}

	// Ограничиваем максимальный период 365 днями
	if days > 365 {
		days = 365
	}

	balance, err := h.analyticsService.PredictBalance(accountID, userID, days)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"predicted_balance": balance,
		"days":              days,
	}); err != nil {
		h.logger.WithError(err).Error("failed to encode response")
	}
}