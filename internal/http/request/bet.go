package request

type PlaceBetRequest struct {
	UserID  string  `json:"user_id"`
	EventID string  `json:"event_id"`
	Odds    float64 `json:"odds"`
	Amount  float64 `json:"amount"`
}

type SettleBetRequest struct {
	EventID string `json:"event_id"`
	Result  string `json:"result"`
}
