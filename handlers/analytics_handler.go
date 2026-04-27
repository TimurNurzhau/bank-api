package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bank-api/middleware"
	"bank-api/services"

	"github.com/gorilla/mux"
)

type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

func (h *AnalyticsHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	stats, err := h.analyticsService.GetMonthlyStats(userID)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	creditLoad, _ := h.analyticsService.GetCreditLoad(userID)

	result := map[string]interface{}{
		"monthly_income":   stats["income"],
		"monthly_expenses": stats["expenses"],
		"credit_load":      creditLoad,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *AnalyticsHandler) PredictBalance(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	vars := mux.Vars(r)

	accountID, err := strconv.Atoi(vars["accountId"])
	if err != nil {
		http.Error(w, `{"error":"invalid account id"}`, http.StatusBadRequest)
		return
	}

	daysStr := r.URL.Query().Get("days")
	days := 30
	if daysStr != "" {
		days, _ = strconv.Atoi(daysStr)
	}

	balance, err := h.analyticsService.PredictBalance(accountID, userID, days)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"predicted_balance": balance,
		"days":              days,
	})
}