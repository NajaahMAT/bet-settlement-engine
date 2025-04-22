package repositories

import (
	"bet-settlement-engine/internal/domain/interface/repo"
	"bet-settlement-engine/internal/domain/model"
	"sync"
)

type betRepositoryImpl struct{}

var (
	bets     = make([]*model.Bet, 0)
	betMutex sync.RWMutex
)

func NewBetRepository() repo.BetRepository {
	return &betRepositoryImpl{}
}

func (r *betRepositoryImpl) SaveBet(bet *model.Bet) {
	betMutex.Lock()
	defer betMutex.Unlock()
	bets = append(bets, bet)
}

func (r *betRepositoryImpl) GetBetsByEvent(eventID string) []*model.Bet {
	betMutex.RLock()
	defer betMutex.RUnlock()
	var results []*model.Bet
	for _, b := range bets {
		if b.EventID == eventID {
			results = append(results, b)
		}
	}
	return results
}
