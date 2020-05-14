package intr

import (
	"time"

	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
)

type timer struct {
	cmp.Component
	duration time.Duration
	left     time.Duration
	state    cmp.TimerState
}

func (t timer) String() string {
	return "{timer:" + t.GetID() + " state:" + t.state.String() + " duration:" + t.duration.String() + "}"
}

func (t *timer) GetState() evt.State {
	event := evt.State{
		EntityID:   t.GetID(),
		State:      t.GetTimerState().String(),
		Attributes: make(map[string]interface{}, 0),
	}

	event.Attributes["duration"] = t.duration
	event.Attributes["left"] = t.left

	return event
}

func (t *timer) GetTimerState() cmp.TimerState {
	return t.state
}

func (t *timer) SetTimerState(state cmp.TimerState) {
	t.state = state
}

func (t *timer) GetDuration() time.Duration {
	return t.duration
}

func (t *timer) GetLeft() time.Duration {
	return t.left
}

func (t *timer) SetLeft(left time.Duration) {
	t.left = left
}

// NewTimer ...
func NewTimer(cfg map[string]string) cmp.Timer {
	d, err := time.ParseDuration(cfg["duration"])
	if err != nil {
		return nil
	}
	return &timer{
		Component: NewComponent(cmp.DomainTimer, cfg["timer"], "internal"),
		duration:  d,
		left:      0,
		state:     cmp.TimerStateInactive,
	}
}
