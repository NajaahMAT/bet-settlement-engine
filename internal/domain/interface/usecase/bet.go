package usecase

type BetUsecase interface {
	PlaceBet(userID, eventID string, odds, amount float64) error
	SettleBet(eventID, result string)
}
