package repo

import "bet-settlement-engine/internal/domain/model"

type UserRepository interface {
	GetUser(userID string) (*model.User, bool)
	CreateUser(userID string) *model.User
	UpdateBalance(userID string, amount float64)
}
