package ptf

import (
	"bytes"
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tsybulin/gosha/cmp"
	cmqtt "github.com/tsybulin/gosha/cmp/mqtt"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/logger"
)

func (mp *mqttPlatform) pushSwitch(c cmp.Component, what string) {
	if c.GetDomain() == cmp.DomainSwitch && (what == "power" || what == "all") {
		if ms, ok := c.(cmqtt.Switch); ok {
			payload := ms.GetPayloadOff()
			if ms.IsOn() {
				payload = ms.GetPayloadOn()
			}

			token := mp.client.Publish(ms.GetCommandTopic(), 0, false, payload)
			token.Wait()

			go func() {
				if mp.eventBus != nil {
					mp.eventBus.Publish(logger.Topic, logger.LevelDebug, "MqttPlatform.Push", c.GetID(), ms.GetCommandTopic(), payload)
				}
			}()
		}
	}
}

func (mp *mqttPlatform) stateHandlerSwitch(c mqttComponent, message mqtt.Message) {
	if c.GetDomain() != cmp.DomainSwitch {
		return
	}

	ms, ok := c.(cmqtt.Switch)

	if !ok {
		return
	}

	if message.Topic() != ms.GetStateTopic() {
		return
	}

	payload := string(message.Payload())

	if ms.GetStateValueTemplate() != nil {
		m := map[string]interface{}{}
		if err := json.Unmarshal(message.Payload(), &m); err != nil {
			return
		}

		var buf bytes.Buffer
		if err := ms.GetStateValueTemplate().Execute(&buf, m); err != nil {
			return
		}

		if buf.Len() == 0 {
			return
		}

		payload = buf.String()
	}

	state := payload == ms.GetPayloadOn()

	if state != ms.IsOn() {
		event := eventFor(c, func() {})

		if state {
			event.Event.Data.NewState.State = "on"
		} else {
			event.Event.Data.NewState.State = "off"
		}

		mp.eventBus.Publish(evt.StateChangedTopic, event)
		mp.eventBus.Publish(logger.Topic, logger.LevelDebug, "MqttPlatform.handleSwitchState power", ms.GetID(), state)
	}
}

func (mp *mqttPlatform) mqttSubscribeSwitch(c cmp.Component) {
	if c.GetDomain() == cmp.DomainSwitch {
		if ms, ok := c.(cmqtt.Switch); ok {
			token := mp.client.Subscribe(ms.GetStateTopic(), 0, mp.stateHandler)
			token.Wait()

			token = mp.client.Publish(ms.GetCommandTopic(), 0, false, " ")
			token.Wait()
		}
	}
}
