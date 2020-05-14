package svc

import (
	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/cmp/mqtt"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/logger"
)

// Mqtt ...
type Mqtt interface {
	Service
}

type mqttService struct {
	Service
	eventBus   evt.Bus
	components map[string]mqtt.Availability
}

func (ms *mqttService) Description() Description {
	return Description{
		ID: ms.GetID(),
	}
}

func (ms *mqttService) registerComponent(c cmp.Component) {
	if c.GetPlatform() != "mqtt" {
		return
	}

	mc, ok := c.(mqtt.Availability)

	if !ok {
		return
	}

	ms.components[c.GetID()] = mc

	go func() {
		ms.eventBus.Publish(logger.Topic, logger.LevelDebug, "MqttService.registerComponent", c.GetID())
	}()
}

func (ms *mqttService) componentStateChanged(event evt.Message) {
	if "event" != event.Type || "state_changed" != event.Event.EventType || event.Event.Data.NewState == nil {
		return
	}

	c := ms.components[event.Event.Data.EntityID]
	if c == nil {
		return
	}

	if event.Event.Data.NewState.Attributes["available"] != nil {
		available, ok := event.Event.Data.NewState.Attributes["available"].(bool)
		if ok && c.GetAvailable() != available {
			c.SetAvailable(available)
			go func() {
				ms.eventBus.Publish(logger.Topic, logger.LevelDebug, "MqttService.componentStateChanged available", c.(cmp.Component).GetID(), c.GetAvailable())
			}()
		}
	}
}

func newMqttService(eventBus evt.Bus) Mqtt {
	ms := &mqttService{
		Service:    newService(cmp.DomainMqtt),
		eventBus:   eventBus,
		components: make(map[string]mqtt.Availability, 0),
	}

	ms.eventBus.Subscribe(cmp.RegisterTopic, "MqttService.registerComponent", ms.registerComponent)
	ms.eventBus.Subscribe(evt.StateChangedTopic, "MqttService.componentStateChanged", ms.componentStateChanged)

	return ms
}
