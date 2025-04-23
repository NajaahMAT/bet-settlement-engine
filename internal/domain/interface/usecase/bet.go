package usecase

import (
	"bet-settlement-engine/internal/http/request"
	"bet-settlement-engine/internal/http/response"

	"github.com/google/uuid"
)

type BetUsecase interface {
	PlaceBet(req request.PlaceBetRequest) (uuid.UUID, error)
	SettleBet(req request.SettleBetRequest) ([]response.SettleBetResponse, error)
}
