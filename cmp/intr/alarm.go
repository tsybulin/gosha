package intr

import (
	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
)

type alarm struct {
	cmp.Component
	state cmp.AlarmState
}

func (a *alarm) AlarmState() cmp.AlarmState {
	return a.state
}

func (a *alarm) SetAlarmState(state cmp.AlarmState) {
	a.state = state
	// log.Println("--- set state", a.GetID(), a.state.String())
}

func (a *alarm) GetState() evt.State {
	return evt.State{
		EntityID:   a.GetID(),
		State:      a.state.String(),
		Attributes: make(map[string]interface{}, 0),
	}
}

// NewAlarm ...
func NewAlarm(cfg map[string]string) cmp.Alarm {
	return &alarm{
		Component: NewComponent(cmp.DomainAlarm, cfg["alarm"], "internal"),
		state:     cmp.AlarmStateDisarmed,
	}
}

// NewAlarmWithPlatform ...
func NewAlarmWithPlatform(cfg map[string]string) cmp.Alarm {
	return &alarm{
		Component: NewComponent(cmp.DomainAlarm, cfg["alarm"], cfg["platform"]),
		state:     cmp.AlarmStateDisarmed,
	}
}
