package mqtt

import (
	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/cmp/intr"
	"github.com/tsybulin/gosha/evt"
)

// Alarm ...
type Alarm interface {
	cmp.Alarm
	Availability
	GetStateTopic() string
	GetCommandTopic() string
}

type alarm struct {
	cmp.Alarm
	Availability
	stateTopic   string
	commandTopic string
}

func (a *alarm) GetStateTopic() string {
	return a.stateTopic
}

func (a *alarm) GetCommandTopic() string {
	return a.commandTopic
}

func (a *alarm) GetState() evt.State {
	state := a.Alarm.GetState()
	state.Attributes["available"] = a.GetAvailable()
	return state
}

// NewMqttAlarm ...
func NewMqttAlarm(cfg map[string]string) Alarm {
	return &alarm{
		Alarm:        intr.NewAlarmWithPlatform(cfg),
		Availability: newMqttAvailability(cfg),
		stateTopic:   cfg["state_topic"],
		commandTopic: cfg["command_topic"],
	}
}
