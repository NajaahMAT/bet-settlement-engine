package response

type PlaceBetResponse struct {
	Msg     string  `json:"msg"`
	BetID   string  `json:"bet_id"`
	UserID  string  `json:"user_id"`
	EventID string  `json:"event_id"`
	Amount  float64 `json:"amount"`
	Odds    float64 `json:"odds"`
	Result  string  `json:"result"`
}

type SettleBetResponse struct {
	Msg       string  `json:"msg"`
	BetID     string  `json:"bet_id"`
	UserID    string  `json:"user_id"`
	AmountWon float64 `json:"amount_won"`
}

type CheckBalanceResponse struct {
	UserID  string  `json:"user_id"`
	Balance float64 `json:"balance"`
}
