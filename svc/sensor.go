package svc

import (
	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/logger"
)

// Sensor ...
type Sensor interface {
	Service
}

type sensorService struct {
	Service
	eventBus evt.Bus
	sensors  map[string]cmp.Sensor
}

func (ss *sensorService) registerSensor(c cmp.Component) {
	if c.GetDomain() != cmp.DomainSensor {
		return
	}

	ss.sensors[c.GetID()] = c.(cmp.Sensor)

	go func() {
		ss.eventBus.Publish(logger.Topic, logger.LevelDebug, "SensorService.registerSensor", c.GetID())
	}()
}

func (ss *sensorService) States() StateResult {
	stateResults := make([]evt.State, 0)

	for _, c := range ss.sensors {
		stateResults = append(stateResults, c.GetState())
	}

	return StateResult{
		ID:      1,
		Type:    "result",
		Success: true,
		Result:  stateResults,
	}
}

func (ss *sensorService) sensorStateChanged(event evt.Message) {
	if "event" != event.Type || "state_changed" != event.Event.EventType || event.Event.Data.NewState == nil {
		return
	}

	s := ss.sensors[event.Event.Data.EntityID]
	if s == nil {
		return
	}

	if s.GetValue() != event.Event.Data.NewState.State {
		s.SetValue(event.Event.Data.NewState.State)
		go func() {
			ss.eventBus.Publish(logger.Topic, logger.LevelDebug, "SensorService.sensorStateChanged value", s.GetID(), s.GetValue())
		}()
	}
}

func newSensorService(eventBus evt.Bus) Sensor {
	ss := &sensorService{
		Service:  newService(cmp.DomainSensor),
		eventBus: eventBus,
		sensors:  make(map[string]cmp.Sensor, 0),
	}

	ss.eventBus.Subscribe(cmp.RegisterTopic, "SensorService.registerSensor", ss.registerSensor)
	ss.eventBus.Subscribe(evt.StateChangedTopic, "SensorService.sensorStateChanged", ss.sensorStateChanged)

	return ss
}
