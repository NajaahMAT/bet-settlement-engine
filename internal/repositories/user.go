package repositories

import (
	"bet-settlement-engine/internal/domain/interface/repo"
	"bet-settlement-engine/internal/domain/model"
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type userRepositoryImpl struct {
	logger      *zap.Logger
	redisClient *redis.Client
}

var userMutexMap = sync.Map{} // map[string]*sync.Mutex

func getUserMutex(userID string) *sync.Mutex {
	mutexIface, _ := userMutexMap.LoadOrStore(userID, &sync.Mutex{})
	return mutexIface.(*sync.Mutex)
}

func NewUserRepository(logger *zap.Logger, redisClient *redis.Client) repo.UserRepository {
	return &userRepositoryImpl{
		logger:      logger,
		redisClient: redisClient,
	}
}

func (r *userRepositoryImpl) GetUser(userID string) (*model.User, bool) {
	start := time.Now()
	val, err := r.redisClient.Get(context.Background(), "user:"+userID).Result()
	if err == redis.Nil {
		r.logger.Warn("User not found in Redis",
			zap.String("user_id", userID),
			zap.Duration("duration", time.Since(start)),
		)
		return nil, false
	} else if err != nil {
		r.logger.Error("Failed to retrieve user from Redis",
			zap.String("user_id", userID),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return nil, false
	}

	var user model.User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		r.logger.Error("Failed to unmarshal user data",
			zap.String("user_id", userID),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return nil, false
	}

	r.logger.Info("User retrieved from Redis",
		zap.String("user_id", userID),
		zap.Float64("balance", user.Balance),
		zap.Duration("duration", time.Since(start)),
	)
	return &user, true
}

func (r *userRepositoryImpl) CreateUser(userID string) *model.User {
	start := time.Now()
	user := &model.User{ID: userID, Balance: 1000.0}
	userData, err := json.Marshal(user)
	if err != nil {
		r.logger.Error("Failed to marshal new user data",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return nil
	}

	err = r.redisClient.Set(context.Background(), "user:"+userID, userData, 0).Err()
	if err != nil {
		r.logger.Error("Failed to store new user in Redis",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return nil
	}

	r.logger.Info("User created in Redis",
		zap.String("user_id", userID),
		zap.Float64("initial_balance", user.Balance),
		zap.Duration("duration", time.Since(start)),
	)
	return user
}

func (r *userRepositoryImpl) UpdateBalance(userID string, amount float64) error {
	start := time.Now()
	mutex := getUserMutex(userID)
	mutex.Lock()
	defer mutex.Unlock()

	val, err := r.redisClient.Get(context.Background(), "user:"+userID).Result()
	if err == redis.Nil {
		r.logger.Warn("Cannot update balance: user not found",
			zap.String("user_id", userID),
			zap.Float64("delta_amount", amount),
			zap.Duration("duration", time.Since(start)),
		)
		return errors.New("user not found")
	} else if err != nil {
		r.logger.Error("Failed to retrieve user for balance update",
			zap.String("user_id", userID),
			zap.Float64("delta_amount", amount),
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		return err
	}

	var user model.User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		r.logger.Error("Failed to unmarshal user during balance update",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return err
	}

	user.Balance += amount
	userData, err := json.Marshal(&user)
	if err != nil {
		r.logger.Error("Failed to marshal updated user data",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return err
	}

	err = r.redisClient.Set(context.Background(), "user:"+userID, userData, 0).Err()
	if err != nil {
		r.logger.Error("Failed to update user balance in Redis",
			zap.String("user_id", userID),
			zap.Float64("new_balance", user.Balance),
			zap.Error(err),
		)
		return err
	}

	r.logger.Info("User balance updated successfully",
		zap.String("user_id", userID),
		zap.Float64("updated_balance", user.Balance),
		zap.Float64("change", amount),
		zap.Duration("duration", time.Since(start)),
	)
	return nil
}
