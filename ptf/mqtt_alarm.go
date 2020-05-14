package ptf

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tsybulin/gosha/cmp"
	cmqtt "github.com/tsybulin/gosha/cmp/mqtt"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/logger"
)

func (mp *mqttPlatform) mqttSubscribeAlarm(c cmp.Component) {
	if c.GetDomain() == cmp.DomainAlarm {
		if al, ok := c.(cmqtt.Alarm); ok {
			token := mp.client.Subscribe(al.GetStateTopic(), 0, mp.stateHandler)
			token.Wait()
		}
	}
}

func (mp *mqttPlatform) pushAlarm(c cmp.Component) {
	if c.GetDomain() == cmp.DomainAlarm {
		if al, ok := c.(cmqtt.Alarm); ok {
			payload := map[cmp.AlarmState]string{
				cmp.AlarmStateDisarmed:   "DISARM",
				cmp.AlarmStateArmedHome:  "ARM_HOME",
				cmp.AlarmStateArmedAway:  "ARM_AWAY",
				cmp.AlarmStateArmedNight: "ARM_NIGHT",
				cmp.AlarmStateTriggered:  "TRIGGER",
			}[al.AlarmState()]

			token := mp.client.Publish(al.GetCommandTopic(), 0, false, payload)
			token.Wait()

			go func() {
				if mp.eventBus != nil {
					mp.eventBus.Publish(logger.Topic, logger.LevelDebug, "MqttPlatform.Push", c.GetID(), al.GetCommandTopic(), payload)
				}
			}()
		}
	}
}

func (mp *mqttPlatform) stateHandlerAlarm(c mqttComponent, message mqtt.Message) {
	if c.GetDomain() != cmp.DomainAlarm {
		return
	}

	al, ok := c.(cmqtt.Alarm)

	if !ok {
		return
	}

	if message.Topic() != al.GetStateTopic() {
		return
	}

	state := map[string]cmp.AlarmState{
		cmp.AlarmStateDisarmed.String():   cmp.AlarmStateDisarmed,
		cmp.AlarmStateArmedHome.String():  cmp.AlarmStateArmedHome,
		cmp.AlarmStateArmedAway.String():  cmp.AlarmStateArmedAway,
		cmp.AlarmStateArmedNight.String(): cmp.AlarmStateArmedNight,
		cmp.AlarmStateTriggered.String():  cmp.AlarmStateTriggered,
	}[string(message.Payload())]

	if state > 0 && state != al.AlarmState() {
		event := eventFor(c, func() {})
		event.Event.Data.NewState.State = state.String()
		mp.eventBus.Publish(evt.StateChangedTopic, event)
		mp.eventBus.Publish(logger.Topic, logger.LevelDebug, "MqttPlatform.handleAlarmState power", al.GetID(), state.String())
	}
}
