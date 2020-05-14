package svc

import (
	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/logger"
)

// Switch ...
type Switch interface {
	Service
	TurnOn(string)
	TurnOff(string)
	Toggle(string)
}

type switchService struct {
	Service
	eventBus evt.Bus
	switches map[string]cmp.Switch
}

func (ss *switchService) States() StateResult {
	stateResults := make([]evt.State, 0)

	for _, c := range ss.switches {
		stateResults = append(stateResults, c.GetState())
	}

	return StateResult{
		ID:      1,
		Type:    "result",
		Success: true,
		Result:  stateResults,
	}
}

func (ss *switchService) Components() []string {
	components := make([]string, 0)
	for id := range ss.switches {
		components = append(components, id)
	}
	return components
}

func (ss *switchService) Description() Description {
	return Description{
		ID: ss.GetID(),
		Methods: []Method{
			{Method: "turn_on", Parameters: []string{"entity_id"}},
			{Method: "turn_off", Parameters: []string{"entity_id"}},
			{Method: "toggle", Parameters: []string{"entity_id"}},
		},
	}
}

func (ss *switchService) registerSwitch(c cmp.Component) {
	if c.GetDomain() != cmp.DomainSwitch {
		return
	}

	sw, ok := c.(cmp.Switch)

	if !ok {
		return
	}

	ss.switches[c.GetID()] = sw

	go func() {
		ss.eventBus.Publish(logger.Topic, logger.LevelDebug, "SwitchService.registerSwitch", c.GetID())
	}()
}

func (ss *switchService) unregisterSwitch(c cmp.Component) {
	sw := ss.switches[c.GetID()]
	if sw == nil {
		return
	}

	delete(ss.switches, c.GetID())
	go func() {
		ss.eventBus.Publish(logger.Topic, logger.LevelDebug, "SwitchService.unregisterSwitch", c.GetID())
	}()
}

func (ss *switchService) switchStateChanged(event evt.Message) {
	if "event" != event.Type || "state_changed" != event.Event.EventType || event.Event.Data.NewState == nil {
		return
	}

	sw := ss.switches[event.Event.Data.EntityID]
	if sw == nil {
		return
	}

	if sw.IsOn() != (event.Event.Data.NewState.State == "on") {
		sw.Toggle()
		go func() {
			ss.eventBus.Publish(logger.Topic, logger.LevelDebug, "SwitchService.switchStateChanged power", sw.GetID(), sw.IsOn())
		}()
	}
}

func (ss *switchService) TurnOn(id string) {
	sw := ss.switches[id]
	if sw == nil {
		ss.eventBus.Publish(logger.Topic, logger.LevelWarn, "SwitchService.TurnOn unknown id:", id)
		return
	}

	if sw.IsOn() {
		return
	}

	ss.eventBus.Publish(evt.StateChangedTopic, eventFor(sw, sw.TurnOn))

	ss.eventBus.Publish(evt.PtfPushTopic, sw, "power")
}

func (ss *switchService) TurnOff(id string) {
	sw := ss.switches[id]
	if sw == nil {
		ss.eventBus.Publish(logger.Topic, logger.LevelWarn, "SwitchService.TurnOff unknown id:", id)
		return
	}

	if !sw.IsOn() {
		return
	}

	ss.eventBus.Publish(evt.StateChangedTopic, eventFor(sw, sw.TurnOff))
	ss.eventBus.Publish(evt.PtfPushTopic, sw, "power")
}

func (ss *switchService) Toggle(id string) {
	sw := ss.switches[id]
	if sw == nil {
		ss.eventBus.Publish(logger.Topic, logger.LevelWarn, "SwitchService.Toggle unknown id:", id)
		return
	}

	ss.eventBus.Publish(evt.StateChangedTopic, eventFor(sw, sw.Toggle))
	ss.eventBus.Publish(evt.PtfPushTopic, sw, "power")
}

func newSwitchService(eventBus evt.Bus) Switch {
	ss := &switchService{
		Service:  newService(cmp.DomainSwitch),
		eventBus: eventBus,
		switches: make(map[string]cmp.Switch, 0),
	}

	ss.eventBus.Subscribe(cmp.RegisterTopic, "SwitchService.registerSwitch", ss.registerSwitch)
	ss.eventBus.Subscribe(evt.StateChangedTopic, "SwitchService.switchStateChanged", ss.switchStateChanged)

	return ss
}
