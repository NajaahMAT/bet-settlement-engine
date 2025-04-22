package usecase

type UserUsecase interface {
	GetOrCreateUser(userID string) float64
	AdjustBalance(userID string, amount float64)
}
