package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"bank-api/middleware"
	"bank-api/models"
	"bank-api/services"

	"github.com/gorilla/mux"
)

type CreditHandler struct {
	creditService *services.CreditService
	cbrService    *services.CBRService
}

func NewCreditHandler(creditService *services.CreditService, cbrService *services.CBRService) *CreditHandler {
	return &CreditHandler{
		creditService: creditService,
		cbrService:    cbrService,
	}
}

func (h *CreditHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req models.CreateCreditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	keyRate, _ := h.cbrService.GetKeyRate()

	credit, err := h.creditService.CreateCredit(userID, &req, keyRate)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(credit); err != nil {
		log.Printf("encode error: %v", err)
	}
}

func (h *CreditHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	credits, err := h.creditService.GetUserCredits(userID)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(credits); err != nil {
		log.Printf("encode error: %v", err)
	}
}

func (h *CreditHandler) GetSchedule(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	vars := mux.Vars(r)
	creditID, err := strconv.Atoi(vars["creditId"])
	if err != nil {
		http.Error(w, `{"error":"invalid credit id"}`, http.StatusBadRequest)
		return
	}

	schedule, err := h.creditService.GetCreditSchedule(creditID, userID)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(schedule); err != nil {
		log.Printf("encode error: %v", err)
	}
}