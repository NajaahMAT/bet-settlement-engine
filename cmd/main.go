package main

import (
	router "bet-settlement-engine/internal/http"
	"bet-settlement-engine/pkg/logger"
	"net/http"

	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()
	defer logger.Log.Sync() // flush logs

	r := router.SetupRoutes()
	logger.Log.Info("Server started", zap.String("port", ":8080"))
	http.ListenAndServe(":8080", r)
}
