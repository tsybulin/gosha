package intr

import (
	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
)

type binarySensor struct {
	cmp.Component
	state       bool
	deviceClass string
}

func (bs *binarySensor) IsOn() bool {
	return bs.state
}

func (bs *binarySensor) SetOn(state bool) {
	bs.state = state
	// log.Println("--- state", bs.GetID(), bs.state)
}

func (bs *binarySensor) GetDeviceClass() string {
	return bs.deviceClass
}

func (bs *binarySensor) GetOnString() string {
	if bs.IsOn() {
		return "on"
	} else {
		return "off"
	}
}

func (bs *binarySensor) GetState() evt.State {
	state := bs.Component.GetState()

	state.Attributes["device_class"] = bs.deviceClass
	state.State = bs.GetOnString()

	return state
}

func newBinarySensor(cfg map[string]string) cmp.BinarySensor {
	return &binarySensor{
		Component:   NewComponent(cmp.DomainBinarySensor, cfg["binary_sensor"], "internal"),
		deviceClass: cfg["device_class"],
		state:       false,
	}
}

func NewBinarySensorWithPlatform(platform string, cfg map[string]string) cmp.BinarySensor {
	return &binarySensor{
		Component:   NewComponent(cmp.DomainBinarySensor, cfg["binary_sensor"], platform),
		deviceClass: cfg["device_class"],
		state:       false,
	}
}
