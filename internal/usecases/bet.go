package usecases

import (
	"bet-settlement-engine/internal/domain/interface/repo"
	"bet-settlement-engine/internal/domain/model"
	"fmt"

	"github.com/google/uuid"
)

type betUsecaseImpl struct {
	betRepo  repo.BetRepository
	userRepo repo.UserRepository
}

func NewBetUsecase(betRepo repo.BetRepository, userRepo repo.UserRepository) *betUsecaseImpl {
	return &betUsecaseImpl{betRepo: betRepo, userRepo: userRepo}
}

func (b *betUsecaseImpl) PlaceBet(userID, eventID string, odds, amount float64) error {
	user, exists := b.userRepo.GetUser(userID)
	if !exists {
		user = b.userRepo.CreateUser(userID)
	}
	if user.Balance < amount {
		return fmt.Errorf("insufficient balance")
	}
	b.userRepo.UpdateBalance(userID, -amount)
	bet := &model.Bet{
		ID:      uuid.New().String(),
		UserID:  userID,
		EventID: eventID,
		Odds:    odds,
		Amount:  amount,
		Result:  "placed",
	}
	b.betRepo.SaveBet(bet)
	return nil
}

func (b *betUsecaseImpl) SettleBet(eventID, result string) {
	bets := b.betRepo.GetBetsByEvent(eventID)
	for _, bet := range bets {
		if bet.Result != "placed" {
			continue
		}
		bet.Result = result
		if result == "win" {
			winnings := bet.Amount * bet.Odds
			b.userRepo.UpdateBalance(bet.UserID, winnings)
		}
	}
}
