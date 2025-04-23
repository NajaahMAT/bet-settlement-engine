package usecase

import (
	"bet-settlement-engine/internal/http/request"

	"github.com/google/uuid"
)

type BetUsecase interface {
	PlaceBet(req request.PlaceBetRequest) (uuid.UUID, error)
	SettleBet(req request.SettleBetRequest)
}
