package mqtt

import "github.com/tsybulin/gosha/cmp"

// Availability ...
type Availability interface {
	cmp.Availability
	AvailabilityTopic() string
	PayloadAvailable() string
	PayloadNotAvailable() string
}

type mqttAvailability struct {
	availabilityTopic   string
	payloadAvailable    string
	payloadNotAvailable string
	available           bool
}

func (mc *mqttAvailability) GetAvailable() bool {
	return mc.available
}

func (mc *mqttAvailability) SetAvailable(available bool) {
	mc.available = available
}

func (mc *mqttAvailability) AvailabilityTopic() string {
	return mc.availabilityTopic
}

func (mc *mqttAvailability) PayloadAvailable() string {
	return mc.payloadAvailable
}

func (mc *mqttAvailability) PayloadNotAvailable() string {
	return mc.payloadNotAvailable
}

func newMqttAvailability(cfg map[string]string) Availability {
	mc := &mqttAvailability{
		availabilityTopic:   cfg["availability_topic"],
		payloadAvailable:    cfg["payload_available"],
		payloadNotAvailable: cfg["payload_not_available"],
		available:           false,
	}

	if len(mc.payloadAvailable) == 0 {
		mc.payloadAvailable = "Online"
	}

	if len(mc.payloadNotAvailable) == 0 {
		mc.payloadNotAvailable = "Offline"
	}

	return mc
}
