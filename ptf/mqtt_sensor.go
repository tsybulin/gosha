package ptf

import (
	"bytes"
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tsybulin/gosha/cmp"
	cmqtt "github.com/tsybulin/gosha/cmp/mqtt"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/logger"
)

func (mp *mqttPlatform) stateHandlerSensor(c mqttComponent, message mqtt.Message) {
	if c.GetDomain() != cmp.DomainSensor {
		return
	}

	ms, ok := c.(cmqtt.Sensor)

	if ok && message.Topic() != ms.GetStateTopic() {
		return
	}

	value := string(message.Payload())

	if ms.GetValueTemplate() != nil {
		m := map[string]interface{}{}
		if err := json.Unmarshal(message.Payload(), &m); err != nil {
			return
		}

		var buf bytes.Buffer
		if err := ms.GetValueTemplate().Execute(&buf, m); err != nil {
			return
		}

		if buf.Len() == 0 {
			return
		}

		value = buf.String()
	}

	if value != ms.GetValue() {
		event := eventFor(c, func() {})
		event.Event.Data.NewState.State = fmt.Sprint(value)

		mp.eventBus.Publish(evt.StateChangedTopic, event)
		mp.eventBus.Publish(logger.Topic, logger.LevelDebug, "MqttPlatform.handleSesorState value", ms.GetID(), value)
	}
}

func (mp *mqttPlatform) mqttSubscribeSensor(c cmp.Component) {
	if c.GetDomain() == cmp.DomainSensor {
		if ms, ok := c.(cmqtt.Sensor); ok {
			token := mp.client.Subscribe(ms.GetStateTopic(), 0, mp.stateHandler)
			token.Wait()
		}
	}
}
