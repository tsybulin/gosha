package ptf

import (
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tsybulin/gosha/cmp"
	cmqtt "github.com/tsybulin/gosha/cmp/mqtt"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/logger"
)

type mqttComponent interface {
	cmp.Component
	cmqtt.Availability
}

type mqttPlatform struct {
	Platform
	clientID   string
	host       string
	login      string
	password   string
	eventBus   evt.Bus
	client     mqtt.Client
	components map[string]mqttComponent
}

func (mp *mqttPlatform) Push(c cmp.Component, what string) {
	if mc := mp.components[c.GetID()]; mc == nil {
		return
	}

	mp.pushSwitch(c, what)
	mp.pushLight(c, what)
	mp.pushAlarm(c)

}

func (mp *mqttPlatform) availableHandler(client mqtt.Client, message mqtt.Message) {
	topic := message.Topic()
	for _, c := range mp.components {
		if c.AvailabilityTopic() != topic {
			continue
		}

		event := eventFor(c, func() {})
		event.Event.Data.NewState.Attributes["available"] = (string(message.Payload()) == c.PayloadAvailable())

		if mp.eventBus != nil {
			mp.eventBus.Publish(evt.StateChangedTopic, event)
			mp.eventBus.Publish(logger.Topic, logger.LevelDebug, "MqttPlatform.availableHandler", c.GetID(), c.GetAvailable())
		}
	}
}

func (mp *mqttPlatform) stateHandler(client mqtt.Client, message mqtt.Message) {
	for _, c := range mp.components {
		mp.stateHandlerSwitch(c, message)
		mp.stateHandlerLight(c, message)
		mp.stateHandlerBinarySensor(c, message)
		mp.stateHandlerSensor(c, message)
		mp.stateHandlerAlarm(c, message)
	}
}

func (mp *mqttPlatform) registerMqttComponent(c cmp.Component) {
	if c.GetPlatform() != "mqtt" {
		return
	}

	mc, ok := c.(mqttComponent)
	if !ok {
		return
	}

	mp.components[c.GetID()] = mc
	token := mp.client.Subscribe(mc.AvailabilityTopic(), 0, mp.availableHandler)
	token.Wait()

	mp.mqttSubscribeSwitch(c)
	mp.mqttSubscribeLight(c)
	mp.mqttSubscribeBinarySensor(c)
	mp.mqttSubscribeSensor(c)
	mp.mqttSubscribeAlarm(c)
}

func (mp *mqttPlatform) unregisterMqttComponent(c cmp.Component) {

}

func (mp *mqttPlatform) Start(eventBus evt.Bus) {
	mp.eventBus = eventBus

	opts := mqtt.NewClientOptions()
	opts.AddBroker(mp.host)
	opts.SetClientID(mp.clientID)
	opts.SetUsername(mp.login)
	opts.SetPassword(mp.password)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	opts.SetDefaultPublishHandler(func(client mqtt.Client, message mqtt.Message) {
		mp.eventBus.Publish(logger.Topic, logger.LevelWarn, "MqttPlatform.defaultHandler", message.Topic(), message.Payload())
	})

	opts.OnConnectionLost = func(client mqtt.Client, e error) {
		mp.eventBus.Publish(logger.Topic, logger.LevelWarn, "MqttPlatform.OnConnectionLost", e.Error())
	}

	opts.OnConnect = func(client mqtt.Client) {
		mp.eventBus.Publish(logger.Topic, logger.LevelInfo, "MqttPlatform.OnConnect")
		for _, c := range mp.components {
			mp.registerMqttComponent(c)
		}

		mp.eventBus.Publish(ReadyTopic, cmp.DomainMqtt)
	}

	mp.client = mqtt.NewClient(opts)
	if token := mp.client.Connect(); token.Wait() && token.Error() != nil {
		mp.eventBus.Publish(logger.Topic, logger.LevelError, "MqttPlatform.Connect", token.Error())
	}

	mp.eventBus.Subscribe(cmp.RegisterTopic, "MqttPlatform.registerMqttComponent", mp.registerMqttComponent)
	mp.eventBus.Subscribe(cmp.UnregisterTopic, "MqttPlatform.unregisterMqttComponent", mp.unregisterMqttComponent)
}

// NewMqttPlatform ...
func NewMqttPlatform(cfg map[string]string) Platform {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "123"
	}
	return &mqttPlatform{
		Platform:   newPlatform("mqtt"),
		clientID:   cfg["mqtt"] + "-" + hostname,
		host:       cfg["host"],
		login:      cfg["login"],
		password:   cfg["password"],
		components: make(map[string]mqttComponent, 0),
	}
}
