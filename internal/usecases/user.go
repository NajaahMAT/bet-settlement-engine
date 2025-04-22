package usecases

import (
	"bet-settlement-engine/internal/domain/interface/repo"
)

type userUsecaseImpl struct {
	userRepo repo.UserRepository
}

func NewUserUsecase(userRepo repo.UserRepository) *userUsecaseImpl {
	return &userUsecaseImpl{userRepo: userRepo}
}

func (u *userUsecaseImpl) GetOrCreateUser(userID string) float64 {
	user, exists := u.userRepo.GetUser(userID)
	if !exists {
		user = u.userRepo.CreateUser(userID)
	}
	return user.Balance
}

func (u *userUsecaseImpl) AdjustBalance(userID string, amount float64) {
	u.userRepo.UpdateBalance(userID, amount)
}
