package ptf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tsybulin/gosha/cmp"
	cmqtt "github.com/tsybulin/gosha/cmp/mqtt"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/logger"
)

func (mp *mqttPlatform) pushLight(c cmp.Component, what string) {
	if c.GetDomain() == cmp.DomainLight {
		if ml, ok := c.(cmqtt.Light); ok {

			// on_command_type: "brightness"
			if ml.IsOn() {
				what = "brightness"
			}

			if what == "power" || what == "all" {
				power := ml.GetPayloadOff()
				if ml.IsOn() {
					power = ml.GetPayloadOn()
				}

				token := mp.client.Publish(ml.GetCommandTopic(), 0, false, power)
				token.Wait()

				go func() {
					if mp.eventBus != nil {
						mp.eventBus.Publish(logger.Topic, logger.LevelDebug, "MqttPlatform.Push", c.GetID(), ml.GetCommandTopic(), power)
					}
				}()
			}

			if what == "brightness" || what == "all" {
				brightness := fmt.Sprintf("%d", ml.GetBrightness())
				token := mp.client.Publish(ml.GetBrightnessCommandTopic(), 0, false, brightness)
				token.Wait()

				go func() {
					if mp.eventBus != nil {
						mp.eventBus.Publish(logger.Topic, logger.LevelDebug, "MqttPlatform.Push", c.GetID(), ml.GetBrightnessCommandTopic(), brightness)
					}
				}()
			}
		}
	}
}

func (mp *mqttPlatform) handleState(ml cmqtt.Light, message mqtt.Message, event *evt.Message) bool {
	changed := false

	if message.Topic() == ml.GetStateTopic() {
		payload := string(message.Payload())

		if ml.GetStateValueTemplate() != nil {
			m := map[string]interface{}{}
			if err := json.Unmarshal(message.Payload(), &m); err != nil {
				return false
			}

			var buf bytes.Buffer
			if err := ml.GetStateValueTemplate().Execute(&buf, m); err != nil {
				return false
			}

			if buf.Len() == 0 {
				return false
			}

			payload = buf.String()
		}

		state := payload == ml.GetPayloadOn()
		if state != ml.IsOn() {
			changed = true
			if state {
				event.Event.Data.NewState.State = "on"
			} else {
				event.Event.Data.NewState.State = "off"
			}
		}
	}

	return changed
}

func (mp *mqttPlatform) handleLight(ml cmqtt.Light, message mqtt.Message, event *evt.Message) bool {
	changed := false

	if message.Topic() == ml.GetBrightnessStateTopic() {
		payload := string(message.Payload())

		if ml.GetBrightnessValueTemplate() != nil {
			m := map[string]interface{}{}
			if err := json.Unmarshal(message.Payload(), &m); err != nil {
				return false
			}

			var buf bytes.Buffer
			if err := ml.GetBrightnessValueTemplate().Execute(&buf, m); err != nil {
				return false
			}

			if buf.Len() == 0 {
				return false
			}

			payload = buf.String()
		}

		if d, err := strconv.Atoi(payload); err == nil {
			brightness := int16(d)
			if ml.GetBrightness() != brightness {
				changed = true
				event.Event.Data.NewState.Attributes["brightness"] = brightness
			}
		}
	}

	return changed
}

func (mp *mqttPlatform) stateHandlerLight(c mqttComponent, message mqtt.Message) {
	if c.GetDomain() != cmp.DomainLight {
		return
	}

	ml, ok := c.(cmqtt.Light)

	if !ok {
		return
	}

	if message.Topic() != ml.GetStateTopic() && message.Topic() != ml.GetBrightnessStateTopic() {
		return
	}

	event := eventFor(c, func() {})
	changed := mp.handleState(ml, message, &event)
	changed2 := mp.handleLight(ml, message, &event)
	changed3 := changed || changed2

	if changed3 {
		mp.eventBus.Publish(evt.StateChangedTopic, event)
		mp.eventBus.Publish(logger.Topic, logger.LevelDebug, "MqttPlatform.stateHandlerLight ", c.GetID())
	}
}

func (mp *mqttPlatform) mqttSubscribeLight(c cmp.Component) {
	if c.GetDomain() != cmp.DomainLight {
		return
	}

	if ml, ok := c.(cmqtt.Light); ok {
		token := mp.client.Subscribe(ml.GetStateTopic(), 0, mp.stateHandler)
		token.Wait()

		token = mp.client.Publish(ml.GetCommandTopic(), 0, false, " ")
		token.Wait()

		if ml.GetStateTopic() != ml.GetBrightnessStateTopic() {
			token = mp.client.Subscribe(ml.GetBrightnessStateTopic(), 0, mp.stateHandler)
			token.Wait()

		}

		token = mp.client.Publish(ml.GetBrightnessCommandTopic(), 0, false, " ")
		token.Wait()
	}
}
