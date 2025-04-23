package repositories

import (
	"bet-settlement-engine/internal/domain/interface/repo"
	"bet-settlement-engine/internal/domain/model"
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type betRepositoryImpl struct {
	logger      *zap.Logger
	redisClient *redis.Client
}

func NewBetRepository(logger *zap.Logger, redisClient *redis.Client) repo.BetRepository {
	return &betRepositoryImpl{
		logger:      logger,
		redisClient: redisClient,
	}
}

func (r *betRepositoryImpl) SaveBet(bet *model.Bet) error {
	start := time.Now()
	betKey := "bets:event:" + bet.EventID

	betData, err := json.Marshal(bet)
	if err != nil {
		r.logger.Error("Failed to marshal bet",
			zap.String("bet_id", bet.ID),
			zap.Error(err),
		)
		return err
	}

	err = r.redisClient.RPush(context.Background(), betKey, betData).Err()
	if err != nil {
		r.logger.Error("Failed to save bet to Redis",
			zap.String("bet_id", bet.ID),
			zap.String("event_id", bet.EventID),
			zap.Error(err),
		)
		return err
	}

	r.logger.Info("Bet saved successfully",
		zap.String("bet_id", bet.ID),
		zap.String("user_id", bet.UserID),
		zap.String("event_id", bet.EventID),
		zap.Float64("amount", bet.Amount),
		zap.Float64("odds", bet.Odds),
		zap.String("result", bet.Result),
		zap.Duration("duration", time.Since(start)),
	)
	return nil
}

func (r *betRepositoryImpl) GetBetsByEvent(eventID string) []*model.Bet {
	start := time.Now()
	betKey := "bets:event:" + eventID
	results := []*model.Bet{}

	data, err := r.redisClient.LRange(context.Background(), betKey, 0, -1).Result()
	if err != nil {
		r.logger.Error("Failed to fetch bets from Redis",
			zap.String("event_id", eventID),
			zap.Error(err),
		)
		return results
	}

	for i, item := range data {
		var b model.Bet
		if err := json.Unmarshal([]byte(item), &b); err != nil {
			r.logger.Warn("Failed to unmarshal bet",
				zap.Int("index", i),
				zap.String("event_id", eventID),
				zap.Error(err),
			)
			continue
		}
		results = append(results, &b)
	}

	r.logger.Info("Fetched bets by event",
		zap.String("event_id", eventID),
		zap.Int("count", len(results)),
		zap.Duration("duration", time.Since(start)),
	)
	return results
}
