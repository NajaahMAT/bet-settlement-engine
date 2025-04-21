package http

import (
	"bet-settlement-engine/internal/http/handler"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/v1/bet/place", handler.PlaceBetHandler).Methods("POST")
	r.HandleFunc("/v1/bet/settle", handler.SettleBetHandler).Methods("PUT")
	r.HandleFunc("/v1/user/balance/{userID}", handler.GetBalanceHandler).Methods("GET")
	return r
}
