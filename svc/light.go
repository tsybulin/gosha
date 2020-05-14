package svc

import (
	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/logger"
)

// Light ...
type Light interface {
	Service
	TurnOn(string)
	TurnOff(string)
	Toggle(string)
	SetBrightness(string, int16)
	SetPower(string, bool, int16)
}

type lightService struct {
	Service
	eventBus evt.Bus
	lights   map[string]cmp.Light
}

func (ls *lightService) States() StateResult {
	stateResults := make([]evt.State, 0)

	for _, c := range ls.lights {
		stateResults = append(stateResults, c.GetState())
	}

	return StateResult{
		ID:      1,
		Type:    "result",
		Success: true,
		Result:  stateResults,
	}
}

func (ls *lightService) Components() []string {
	components := make([]string, 0)
	for id := range ls.lights {
		components = append(components, id)
	}
	return components
}

func (ls *lightService) Description() Description {
	return Description{
		ID: ls.GetID(),
		Methods: []Method{
			{Method: "turn_on", Parameters: []string{"entity_id"}},
			{Method: "turn_off", Parameters: []string{"entity_id"}},
			{Method: "toggle", Parameters: []string{"entity_id"}},
			{Method: "10", Parameters: []string{"entity_id"}},
			{Method: "20", Parameters: []string{"entity_id"}},
			{Method: "30", Parameters: []string{"entity_id"}},
			{Method: "40", Parameters: []string{"entity_id"}},
			{Method: "50", Parameters: []string{"entity_id"}},
			{Method: "60", Parameters: []string{"entity_id"}},
			{Method: "70", Parameters: []string{"entity_id"}},
			{Method: "80", Parameters: []string{"entity_id"}},
			{Method: "90", Parameters: []string{"entity_id"}},
			{Method: "100", Parameters: []string{"entity_id"}},
		},
	}
}

func (ls *lightService) registerLight(c cmp.Component) {
	if c.GetDomain() != cmp.DomainLight {
		return
	}

	l, ok := c.(cmp.Light)

	if !ok {
		return
	}

	ls.lights[c.GetID()] = l

	go func() {
		ls.eventBus.Publish(logger.Topic, logger.LevelDebug, "LightService.registerLight", c.GetID())
	}()
}

func (ls *lightService) lightStateChanged(event evt.Message) {
	if "event" != event.Type || "state_changed" != event.Event.EventType || event.Event.Data.NewState == nil {
		return
	}

	l := ls.lights[event.Event.Data.EntityID]
	if l == nil {
		return
	}

	if l.IsOn() != (event.Event.Data.NewState.State == "on") {
		l.Toggle()
		go func() {
			ls.eventBus.Publish(logger.Topic, logger.LevelDebug, "LightService.lightStateChanged power", l.GetID(), l.IsOn())
		}()
	}

	if event.Event.Data.NewState.Attributes["brightness"] != nil {
		brightness, ok := event.Event.Data.NewState.Attributes["brightness"].(int16)
		if ok && l.GetBrightness() != brightness {
			l.SetBrightness(brightness)
			go func() {
				ls.eventBus.Publish(logger.Topic, logger.LevelDebug, "LightService.lightStateChanged brightness", l.GetID(), l.GetBrightness())
			}()
		}
	}

}

func (ls *lightService) TurnOn(id string) {
	l := ls.lights[id]
	if l == nil {
		ls.eventBus.Publish(logger.Topic, logger.LevelWarn, "LightService.TurnOn unknown id:", id)
		return
	}

	if l.IsOn() {
		return
	}

	ls.eventBus.Publish(evt.StateChangedTopic, eventFor(l, l.TurnOn))
	ls.eventBus.Publish(evt.PtfPushTopic, l, "power")
}

func (ls *lightService) TurnOff(id string) {
	l := ls.lights[id]
	if l == nil {
		ls.eventBus.Publish(logger.Topic, logger.LevelWarn, "LightService.TurnOff unknown id:", id)
		return
	}

	if !l.IsOn() {
		return
	}

	ls.eventBus.Publish(evt.StateChangedTopic, eventFor(l, l.TurnOff))
	ls.eventBus.Publish(evt.PtfPushTopic, l, "power")
}

func (ls *lightService) Toggle(id string) {
	l := ls.lights[id]
	if l == nil {
		ls.eventBus.Publish(logger.Topic, logger.LevelWarn, "LightService.Toggle unknown id:", id)
		return
	}

	ls.eventBus.Publish(evt.StateChangedTopic, eventFor(l, l.Toggle))
	ls.eventBus.Publish(evt.PtfPushTopic, l, "power")
}

func (ls *lightService) SetBrightness(id string, brightness int16) {
	l := ls.lights[id]
	if l == nil {
		ls.eventBus.Publish(logger.Topic, logger.LevelWarn, "LightService.SetBrightness unknown id:", id)
		return
	}

	if l.GetBrightness() == brightness {
		return
	}

	ls.eventBus.Publish(evt.StateChangedTopic, eventFor(l, func() {
		l.SetBrightness(brightness)
	}))

	ls.eventBus.Publish(evt.PtfPushTopic, l, "brightness")
}

func (ls *lightService) SetPower(id string, power bool, brightness int16) {
	l := ls.lights[id]
	if l == nil {
		ls.eventBus.Publish(logger.Topic, logger.LevelWarn, "LightService.SetPower unknown id:", id)
		return
	}

	ls.eventBus.Publish(evt.StateChangedTopic, eventFor(l, func() {
		if power {
			l.TurnOn()
		} else {
			l.TurnOff()
		}
		l.SetBrightness(brightness)
	}))

	ls.eventBus.Publish(evt.PtfPushTopic, l, "all")
}

func newLightService(eventBus evt.Bus) Light {
	ls := &lightService{
		Service:  newService(cmp.DomainLight),
		eventBus: eventBus,
		lights:   make(map[string]cmp.Light, 0),
	}

	ls.eventBus.Subscribe(cmp.RegisterTopic, "LightService.registerLight", ls.registerLight)
	ls.eventBus.Subscribe(evt.StateChangedTopic, "LightService.lightStateChanged", ls.lightStateChanged)

	return ls
}
