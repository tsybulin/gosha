package cmp

import "time"

type TimerState int

const (
	TimerStateInactive TimerState = iota
	TimerStateActive
)

func (s TimerState) String() string {
	return [...]string{"inactive", "active"}[s]
}

// Timer ...
type Timer interface {
	Component
	GetTimerState() TimerState
	SetTimerState(TimerState)
	GetDuration() time.Duration
	GetLeft() time.Duration
	SetLeft(time.Duration)
}
