package svc

import (
	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/logger"
)

// Alarm ...
type Alarm interface {
	Service
	SetState(string, string)
	Disarm(string)
	ArmHome(string)
	ArmAway(string)
	ArmNight(string)
	Trigger(string)
}

type alarm struct {
	Service
	eventBus evt.Bus
	alarms   map[string]cmp.Alarm
}

func (a *alarm) States() StateResult {
	stateResults := make([]evt.State, 0)

	for _, c := range a.alarms {
		stateResults = append(stateResults, c.GetState())
	}

	return StateResult{
		ID:      1,
		Type:    "result",
		Success: true,
		Result:  stateResults,
	}
}

func (a *alarm) Components() []string {
	components := make([]string, 0)
	for id := range a.alarms {
		components = append(components, id)
	}
	return components
}

func (a *alarm) Description() Description {
	return Description{
		ID: a.GetID(),
		Methods: []Method{
			{Method: "disarm", Parameters: []string{"entity_id"}},
			{Method: "arm_home", Parameters: []string{"entity_id"}},
			{Method: "arm_away", Parameters: []string{"entity_id"}},
			{Method: "arm_night", Parameters: []string{"entity_id"}},
			{Method: "trigger", Parameters: []string{"entity_id"}},
		},
	}
}

func (a *alarm) registerAlarm(c cmp.Component) {
	if c.GetDomain() != cmp.DomainAlarm {
		return
	}

	al, ok := c.(cmp.Alarm)

	if !ok {
		return
	}

	a.alarms[c.GetID()] = al

	go func() {
		a.eventBus.Publish(logger.Topic, logger.LevelDebug, "AlarmService.registerAlarm", c.GetID())
	}()
}

func (a *alarm) setAlarmState(al cmp.Alarm, state cmp.AlarmState) {
	if al.AlarmState() == state {
		return
	}

	go func() {
		a.eventBus.Publish(evt.StateChangedTopic, eventFor(al, func() {
			al.SetAlarmState(state)
		}))

		a.eventBus.Publish(evt.PtfPushTopic, al, "state")
	}()
}

func (a *alarm) Disarm(id string) {
	al := a.alarms[id]
	if al == nil {
		go func() {
			a.eventBus.Publish(logger.Topic, logger.LevelWarn, "AlarmService.Disarm unknown id", id)
		}()
		return
	}

	a.setAlarmState(al, cmp.AlarmStateDisarmed)
}

func (a *alarm) ArmHome(id string) {
	al := a.alarms[id]
	if al == nil {
		go func() {
			a.eventBus.Publish(logger.Topic, logger.LevelWarn, "AlarmService.ArmHome unknown id", id)
		}()
		return
	}

	a.setAlarmState(al, cmp.AlarmStateArmedHome)
}

func (a *alarm) ArmAway(id string) {
	al := a.alarms[id]
	if al == nil {
		go func() {
			a.eventBus.Publish(logger.Topic, logger.LevelWarn, "AlarmService.ArmAway unknown id", id)
		}()
		return
	}

	a.setAlarmState(al, cmp.AlarmStateArmedAway)
}

func (a *alarm) ArmNight(id string) {
	al := a.alarms[id]
	if al == nil {
		go func() {
			a.eventBus.Publish(logger.Topic, logger.LevelWarn, "AlarmService.ArmNight unknown id", id)
		}()
		return
	}

	a.setAlarmState(al, cmp.AlarmStateArmedNight)
}

func (a *alarm) Trigger(id string) {
	al := a.alarms[id]
	if al == nil {
		go func() {
			a.eventBus.Publish(logger.Topic, logger.LevelWarn, "AlarmService.Trigger unknown id", id)
		}()
		return
	}

	a.setAlarmState(al, cmp.AlarmStateTriggered)
}

func (a *alarm) SetState(id, state string) {
	al := a.alarms[id]
	if al == nil {
		go func() {
			a.eventBus.Publish(logger.Topic, logger.LevelWarn, "AlarmService.SetState unknown id", id)
		}()
		return
	}

	if al.AlarmState().String() != (state) {
		changed := true
		switch state {
		case cmp.AlarmStateDisarmed.String():
			a.Disarm(id)
		case cmp.AlarmStateArmedHome.String():
			a.ArmHome(id)
		case cmp.AlarmStateArmedAway.String():
			a.ArmAway(id)
		case cmp.AlarmStateArmedNight.String():
			a.ArmNight(id)
		case cmp.AlarmStateTriggered.String():
			a.Trigger(id)
		default:
			changed = false
			go func() {
				a.eventBus.Publish(logger.Topic, logger.LevelWarn, "AlarmService.SetState unknown state", state)
			}()
		}
		if changed {
			go func() {
				a.eventBus.Publish(logger.Topic, logger.LevelDebug, "Alarm.SetState state", al.GetID(), al.AlarmState().String())
			}()
		}
	}

}

func (a *alarm) alarmStateChanged(event evt.Message) {
	if "event" != event.Type || "state_changed" != event.Event.EventType || event.Event.Data.NewState == nil {
		return
	}

	al := a.alarms[event.Event.Data.EntityID]
	if al == nil {
		return
	}

	a.SetState(event.Event.Data.EntityID, event.Event.Data.NewState.State)
}

func newAlarmService(eventBus evt.Bus) Alarm {
	a := &alarm{
		Service:  newService(cmp.DomainAlarm),
		eventBus: eventBus,
		alarms:   make(map[string]cmp.Alarm, 0),
	}

	a.eventBus.Subscribe(cmp.RegisterTopic, "AlarmService.registerAlarm", a.registerAlarm)
	a.eventBus.Subscribe(evt.StateChangedTopic, "AlarmService.alarmStateChanged", a.alarmStateChanged)

	return a
}
