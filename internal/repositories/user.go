package repositories

import (
	"bet-settlement-engine/internal/domain/interface/repo"
	"bet-settlement-engine/internal/domain/model"
	"sync"
)

type userRepositoryImpl struct{}

var (
	users     = make(map[string]*model.User)
	userMutex sync.RWMutex
)

func NewUserRepository() repo.UserRepository {
	return &userRepositoryImpl{}
}

func (r *userRepositoryImpl) GetUser(userID string) (*model.User, bool) {
	userMutex.RLock()
	defer userMutex.RUnlock()
	user, exists := users[userID]
	return user, exists
}

func (r *userRepositoryImpl) CreateUser(userID string) *model.User {
	userMutex.Lock()
	defer userMutex.Unlock()
	user := &model.User{ID: userID, Balance: 1000.0}
	users[userID] = user
	return user
}

func (r *userRepositoryImpl) UpdateBalance(userID string, amount float64) {
	userMutex.Lock()
	defer userMutex.Unlock()
	users[userID].Balance += amount
}
