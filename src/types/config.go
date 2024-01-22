package types

import "time"

type ActionLimitConfig struct {
	ActionType string
	Limit      int
	Duration   time.Duration
}

type AllowedCombination struct {
	UserID string
	IP     string
}
