package ptf

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tsybulin/gosha/cmp"
	cmqtt "github.com/tsybulin/gosha/cmp/mqtt"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/logger"
)

func (mp *mqttPlatform) stateHandlerBinarySensor(c mqttComponent, message mqtt.Message) {
	if c.GetDomain() != cmp.DomainBinarySensor {
		return
	}

	bs, ok := c.(cmqtt.BinarySensor)

	if ok && message.Topic() != bs.GetStateTopic() {
		return
	}

	state := string(message.Payload()) == bs.GetPayloadOn()

	if state != bs.IsOn() {
		event := eventFor(c, func() {})
		if state {
			event.Event.Data.NewState.State = "on"
		} else {
			event.Event.Data.NewState.State = "off"
		}

		mp.eventBus.Publish(evt.StateChangedTopic, event)
		mp.eventBus.Publish(logger.Topic, logger.LevelDebug, "MqttPlatform.handleBinarySesorState state", bs.GetID(), state)
	}
}

func (mp *mqttPlatform) mqttSubscribeBinarySensor(c cmp.Component) {
	if c.GetDomain() == cmp.DomainBinarySensor {
		if bs, ok := c.(cmqtt.BinarySensor); ok {
			token := mp.client.Subscribe(bs.GetStateTopic(), 0, mp.stateHandler)
			token.Wait()
		}
	}
}
