package usecases

import (
	"bet-settlement-engine/internal/domain/interface/repo"

	"go.uber.org/zap"
)

type userUsecaseImpl struct {
	userRepo repo.UserRepository
	logger   *zap.Logger
}

func NewUserUsecase(userRepo repo.UserRepository, logger *zap.Logger) *userUsecaseImpl {
	return &userUsecaseImpl{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (u *userUsecaseImpl) GetOrCreateUser(userID string) float64 {
	u.logger.Info("Fetching user", zap.String("user_id", userID))
	user, exists := u.userRepo.GetUser(userID)
	if !exists {
		u.logger.Warn("User not found, creating new user", zap.String("user_id", userID))

		user = u.userRepo.CreateUser(userID)
	} else {
		u.logger.Debug("User found", zap.String("user_id", userID), zap.Float64("balance", user.Balance))
	}

	u.logger.Info("User balance fetched", zap.String("user_id", userID), zap.Float64("balance", user.Balance))

	return user.Balance
}

func (u *userUsecaseImpl) AdjustBalance(userID string, amount float64) {
	u.logger.Info("Adjusting user balance",
		zap.String("user_id", userID),
		zap.Float64("adjustment_amount", amount),
	)

	err := u.userRepo.UpdateBalance(userID, amount)
	if err != nil {
		u.logger.Error("Failed to adjust user balance",
			zap.String("user_id", userID),
			zap.Float64("amount", amount),
			zap.Error(err),
		)
	} else {
		u.logger.Info("User balance adjusted successfully",
			zap.String("user_id", userID),
			zap.Float64("amount", amount),
		)
	}
}
