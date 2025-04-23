package usecases

import (
	"bet-settlement-engine/internal/domain/interface/repo"
	"bet-settlement-engine/internal/domain/model"
	"bet-settlement-engine/internal/http/request"
	"bet-settlement-engine/internal/http/response"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type betUsecaseImpl struct {
	betRepo  repo.BetRepository
	userRepo repo.UserRepository
	logger   *zap.Logger
}

func NewBetUsecase(betRepo repo.BetRepository, userRepo repo.UserRepository, logger *zap.Logger) *betUsecaseImpl {
	return &betUsecaseImpl{
		betRepo:  betRepo,
		userRepo: userRepo,
		logger:   logger,
	}
}

func (b *betUsecaseImpl) PlaceBet(req request.PlaceBetRequest) (res uuid.UUID, err error) {
	b.logger.Info("Placing bet",
		zap.String("user_id", req.UserID),
		zap.Float64("amount", req.Amount),
		zap.String("event_id", req.EventID),
		zap.Float64("odds", req.Odds),
	)
	user, exists := b.userRepo.GetUser(req.UserID)
	if !exists {
		b.logger.Warn("User not found, creating new user", zap.String("user_id", req.UserID))
		user = b.userRepo.CreateUser(req.UserID)
	} else {
		b.logger.Debug("User found", zap.String("user_id", user.ID), zap.Float64("balance", user.Balance))
	}

	if user.Balance < req.Amount {
		b.logger.Warn("Insufficient balance", zap.String("user_id", req.UserID), zap.Float64("balance", user.Balance), zap.Float64("attempted_bet", req.Amount))
		return res, fmt.Errorf("insufficient balance")
	}

	err = b.userRepo.UpdateBalance(req.UserID, -req.Amount)
	if err != nil {
		b.logger.Error("Failed to deduct balance", zap.String("user_id", req.UserID), zap.Error(err))
		return res, err
	}
	b.logger.Info("Balance deducted", zap.String("user_id", req.UserID), zap.Float64("amount", req.Amount))

	betID := uuid.New()
	bet := &model.Bet{
		ID:      betID.String(),
		UserID:  req.UserID,
		EventID: req.EventID,
		Odds:    req.Odds,
		Amount:  req.Amount,
		Result:  "placed",
	}
	err = b.betRepo.SaveBet(bet)
	if err != nil {
		b.logger.Error("Failed to save bet", zap.Any("bet", bet), zap.Error(err))
		return res, err
	}

	b.logger.Info("Bet successfully placed", zap.String("bet_id", bet.ID), zap.String("user_id", req.UserID))

	return betID, nil
}

func (b *betUsecaseImpl) SettleBet(req request.SettleBetRequest) ([]response.SettleBetResponse, error) {
	b.logger.Info("Settling bets", zap.String("event_id", req.EventID), zap.String("result", req.Result))

	bets := b.betRepo.GetBetsByEvent(req.EventID)
	b.logger.Debug("Fetched bets for event", zap.Int("count", len(bets)), zap.String("event_id", req.EventID))

	var responses []response.SettleBetResponse
	var failedUpdates int

	for _, bet := range bets {
		if bet.Result != "placed" {
			b.logger.Debug("Skipping already settled bet", zap.String("bet_id", bet.ID), zap.String("current_result", bet.Result))
			continue
		}

		bet.Result = req.Result
		b.logger.Info("Updating bet result", zap.String("bet_id", bet.ID), zap.String("new_result", req.Result))

		var amountWon float64

		if bet.Result == "win" {
			amountWon = bet.Amount * bet.Odds
			err := b.userRepo.UpdateBalance(bet.UserID, amountWon)
			if err != nil {
				failedUpdates++
				b.logger.Error("Failed to update user balance after win",
					zap.String("user_id", bet.UserID),
					zap.Float64("winnings", amountWon),
					zap.Error(err),
				)
				continue
			} else {
				b.logger.Info("User winnings credited", zap.String("user_id", bet.UserID), zap.Float64("amount", amountWon))
			}
		}

		responses = append(responses, response.SettleBetResponse{
			Msg:       "Bet settled successfully",
			BetID:     bet.ID,
			UserID:    bet.UserID,
			AmountWon: amountWon,
		})
	}

	if len(responses) == 0 {
		return nil, fmt.Errorf("no bets could be settled for event: %s", req.EventID)
	}

	if failedUpdates > 0 {
		b.logger.Warn("Some bets failed to settle completely", zap.Int("failed", failedUpdates))
	}

	return responses, nil
}
