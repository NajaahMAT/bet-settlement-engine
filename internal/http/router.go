package http

import (
	"bet-settlement-engine/internal/http/handler"
	"bet-settlement-engine/internal/repositories"
	"bet-settlement-engine/internal/usecases"
	"bet-settlement-engine/pkg/db"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func SetupRoutes() *mux.Router {

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	// Initialize Redis
	redisDB := &db.RedisDB{}
	if err := redisDB.Init(); err != nil {
		logger.Fatal("Failed to initialize Redis", zap.Error(err))
	}

	redisClient := redisDB.GetClient()

	userRepo := repositories.NewUserRepository(logger, redisClient)
	betRepo := repositories.NewBetRepository(logger, redisClient)

	betUsecase := usecases.NewBetUsecase(betRepo, userRepo, logger)
	userUsecase := usecases.NewUserUsecase(userRepo, logger)

	h := handler.NewHandler(betUsecase, userUsecase)

	r := mux.NewRouter()
	r.HandleFunc("/v1/bet/place", h.PlaceBetHandler).Methods("POST")
	r.HandleFunc("/v1/bet/settle", h.SettleBetHandler).Methods("PUT")
	r.HandleFunc("/v1/user/balance/{userID}", h.GetBalanceHandler).Methods("GET")

	return r
}
