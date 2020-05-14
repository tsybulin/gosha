package mqtt

import (
	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/cmp/intr"
	"github.com/tsybulin/gosha/evt"
)

// BinarySensor ...
type BinarySensor interface {
	cmp.BinarySensor
	Availability
	GetStateTopic() string
	GetPayloadOn() string
	GetPayloadOff() string
}

type mqttBinarySensor struct {
	cmp.BinarySensor
	Availability
	stateTopic string
	payloadOn  string
	payloadOff string
}

func (ms *mqttBinarySensor) GetStateTopic() string {
	return ms.stateTopic
}

func (ms *mqttBinarySensor) GetPayloadOn() string {
	return ms.payloadOn
}

func (ms *mqttBinarySensor) GetPayloadOff() string {
	return ms.payloadOff
}

func (ms *mqttBinarySensor) GetState() evt.State {
	state := ms.BinarySensor.GetState()
	state.Attributes["available"] = ms.GetAvailable()
	return state
}

func NewMqttBinarySensor(cfg map[string]string) BinarySensor {
	ms := &mqttBinarySensor{
		BinarySensor: intr.NewBinarySensorWithPlatform("mqtt", cfg),
		Availability: newMqttAvailability(cfg),
		stateTopic:   cfg["state_topic"],
		payloadOn:    cfg["payload_on"],
		payloadOff:   cfg["payload_off"],
	}

	return ms
}
