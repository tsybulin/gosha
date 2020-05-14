package intr

import (
	"fmt"

	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
)

type sensor struct {
	cmp.Component
	value             interface{}
	deviceClass       string
	unitOfMeasurement string
}

func (s *sensor) GetValue() interface{} {
	return s.value
}

func (s *sensor) SetValue(v interface{}) {
	s.value = v
	// log.Println("--- value", s.GetID(), s.value)
}

func (s *sensor) GetDeviceClass() string {
	return s.deviceClass
}

func (s *sensor) GetUnitOfMeasurement() string {
	return s.unitOfMeasurement
}

func (s *sensor) GetState() evt.State {
	state := s.Component.GetState()

	state.Attributes["device_class"] = s.deviceClass
	state.State = fmt.Sprint(s.value)
	state.Attributes["unit_of_measurement"] = s.unitOfMeasurement

	return state
}

func newSensor(cfg map[string]string) cmp.Sensor {
	return &sensor{
		Component:         NewComponent(cmp.DomainSensor, cfg["sensor"], "internal"),
		deviceClass:       cfg["device_class"],
		unitOfMeasurement: cfg["unit_of_measurement"],
	}
}

func NewSensorWithPlatform(platform string, cfg map[string]string) cmp.Sensor {
	return &sensor{
		Component:         NewComponent(cmp.DomainSensor, cfg["sensor"], platform),
		deviceClass:       cfg["device_class"],
		unitOfMeasurement: cfg["unit_of_measurement"],
	}
}
