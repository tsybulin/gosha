package cmp

import "time"

// Countdown ...
type Countdown interface {
	Component
	SetDuration(time.Duration)
	GetDuration() time.Duration
	GetLeft() time.Duration
	SetLeft(time.Duration)
}
