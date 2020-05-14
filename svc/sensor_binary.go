package svc

import (
	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/logger"
)

// BinarySensor ...
type BinarySensor interface {
	Service
}

type binarySensorService struct {
	Service
	eventBus evt.Bus
	sensors  map[string]cmp.BinarySensor
}

func (ss *binarySensorService) States() StateResult {
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

func (ss *binarySensorService) registerBinarySensor(c cmp.Component) {
	if c.GetDomain() != cmp.DomainBinarySensor {
		return
	}

	ss.sensors[c.GetID()] = c.(cmp.BinarySensor)

	go func() {
		ss.eventBus.Publish(logger.Topic, logger.LevelDebug, "SensorService.registerSensor", c.GetID())
	}()
}

func (ss *binarySensorService) binarySensorStateChanged(event evt.Message) {
	if "event" != event.Type || "state_changed" != event.Event.EventType || event.Event.Data.NewState == nil {
		return
	}

	bs := ss.sensors[event.Event.Data.EntityID]
	if bs == nil {
		return
	}

	if bs.IsOn() != (event.Event.Data.NewState.State == "on") {
		bs.SetOn(event.Event.Data.NewState.State == "on")
		go func() {
			ss.eventBus.Publish(logger.Topic, logger.LevelDebug, "SensorService.sensorStateChanged power", bs.GetID(), bs.IsOn())
		}()
	}
}

func newBinarySensorService(eventBus evt.Bus) BinarySensor {
	ss := &binarySensorService{
		Service:  newService(cmp.DomainBinarySensor),
		eventBus: eventBus,
		sensors:  make(map[string]cmp.BinarySensor, 0),
	}

	ss.eventBus.Subscribe(cmp.RegisterTopic, "BinarySensorService.registerBinarySensor", ss.registerBinarySensor)
	ss.eventBus.Subscribe(evt.StateChangedTopic, "BinarySensorService.binarySensorStateChanged", ss.binarySensorStateChanged)

	return ss
}
