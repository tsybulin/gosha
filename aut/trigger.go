package aut

import (
	"time"

	"github.com/tsybulin/gosha/evt"
)

// Trigger ...
type Trigger interface {
	GetPlatform() string
	Fire() bool
}

// StateTrigger ...
type StateTrigger interface {
	Trigger
	FireState(evt.Event) bool
	GetEntityID() string
	GetFrom() string
	GetTo() string
}

// TimeTrigger ...
type TimeTrigger interface {
	Trigger
	GetAt() time.Time
	FireTime(time.Time) bool
}
