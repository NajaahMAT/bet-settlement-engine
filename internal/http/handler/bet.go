package handler

import (
	"bet-settlement-engine/internal/domain/interface/usecase"
	"bet-settlement-engine/internal/http/request"
	"bet-settlement-engine/internal/http/response"
	"bet-settlement-engine/pkg/logger"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Handler struct {
	betUsecase  usecase.BetUsecase
	userUsecase usecase.UserUsecase
}

func NewHandler(betUC usecase.BetUsecase, userUC usecase.UserUsecase) *Handler {
	return &Handler{
		betUsecase:  betUC,
		userUsecase: userUC,
	}
}

func (h *Handler) PlaceBetHandler(w http.ResponseWriter, r *http.Request) {
	var req request.PlaceBetRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	logger.Log.Info("Placing bet", zap.String("user_id", req.UserID), zap.Float64("amount", req.Amount))

	betID, err := h.betUsecase.PlaceBet(req)
	if err != nil {
		logger.Log.Error("Failed to place bet", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := response.PlaceBetResponse{
		Msg:     "Bet saved successfully",
		BetID:   betID.String(),
		UserID:  req.UserID,
		EventID: req.EventID,
		Amount:  req.Amount,
		Odds:    req.Odds,
		Result:  "placed",
	}

	logger.Log.Info("Bet placed successfully", zap.String("user_id", req.UserID))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) SettleBetHandler(w http.ResponseWriter, r *http.Request) {
	var req request.SettleBetRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	logger.Log.Info("Settling bets", zap.String("event_id", req.EventID), zap.String("result", req.Result))
	h.betUsecase.SettleBet(req)

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	balance := h.userUsecase.GetOrCreateUser(userID)

	logger.Log.Info("Fetched balance", zap.String("user_id", userID), zap.Float64("balance", balance))

	resp := map[string]float64{
		"balance": balance,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
