package svc

import (
	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/logger"
)

// SensorWeather ...
type SensorWeather interface {
	Service
}

type sensorWeather struct {
	Service
	eventBus evt.Bus
	sensors  map[string]cmp.WeatherSensor
}

func (sw *sensorWeather) registerSensor(c cmp.Component) {
	if c.GetDomain() != cmp.DomainWeather {
		return
	}

	sw.sensors[c.GetID()] = c.(cmp.WeatherSensor)

	go func() {
		sw.eventBus.Publish(logger.Topic, logger.LevelDebug, "WeatherService.registerSensor", c.GetID())
	}()
}

func (sw *sensorWeather) sensorStateChanged(event evt.Message) {
	if "event" != event.Type || "state_changed" != event.Event.EventType || event.Event.Data.NewState == nil {
		return
	}

	s := sw.sensors[event.Event.Data.EntityID]
	if s == nil {
		return
	}

	s.SetValues(event)

	go func() {
		sw.eventBus.Publish(logger.Topic, logger.LevelDebug, "WeatherService.sensorStateChanged values", s.GetID())
	}()
}

func (sw *sensorWeather) States() StateResult {
	stateResults := make([]evt.State, 0)

	for _, c := range sw.sensors {
		stateResults = append(stateResults, c.GetState())
	}

	return StateResult{
		ID:      1,
		Type:    "result",
		Success: true,
		Result:  stateResults,
	}
}

func newSensorWeather(eventBus evt.Bus) SensorWeather {
	sw := &sensorWeather{
		Service:  newService(cmp.DomainSensor),
		eventBus: eventBus,
		sensors:  make(map[string]cmp.WeatherSensor, 0),
	}

	sw.eventBus.Subscribe(cmp.RegisterTopic, "SensorWeather.registerSensor", sw.registerSensor)
	sw.eventBus.Subscribe(evt.StateChangedTopic, "SensorWeather.sensorStateChanged", sw.sensorStateChanged)

	return sw
}
