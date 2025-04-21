package model

type Bet struct {
	ID      string
	UserID  string
	EventID string
	Odds    float64
	Amount  float64
	Result  string // placed,win,lose
}
