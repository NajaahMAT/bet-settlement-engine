package repo

import "bet-settlement-engine/internal/domain/model"

type BetRepository interface {
	SaveBet(bet *model.Bet) error
	GetBetsByEvent(eventID string) []*model.Bet
}
