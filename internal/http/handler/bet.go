package handler

import (
	"fmt"
	"net/http"
)

func PlaceBetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Place Bet")
}

func SettleBetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Settle Bet")
}

func GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Bet Balance")
}
