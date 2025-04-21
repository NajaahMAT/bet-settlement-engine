package main

import (
	router "bet-settlement-engine/internal/http"
	"log"
	"net/http"
)

func main() {
	r := router.SetupRoutes()
	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", r)
}
