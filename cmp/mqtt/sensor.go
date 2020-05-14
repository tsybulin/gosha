package mqtt

import (
	"html/template"

	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/cmp/intr"
	"github.com/tsybulin/gosha/evt"
)

// Sensor ...
type Sensor interface {
	cmp.Sensor
	Availability
	GetStateTopic() string
	GetValueTemplate() *template.Template
}

type mqttSensor struct {
	cmp.Sensor
	Availability
	stateTopic    string
	valueTemplate *template.Template
}

func (ms *mqttSensor) GetStateTopic() string {
	return ms.stateTopic
}

func (ms *mqttSensor) GetValueTemplate() *template.Template {
	return ms.valueTemplate
}

func (ms *mqttSensor) GetState() evt.State {
	state := ms.Sensor.GetState()
	state.Attributes["available"] = ms.GetAvailable()
	return state
}

// NewMqttSensor ...
func NewMqttSensor(cfg map[string]string) Sensor {
	ms := &mqttSensor{
		Sensor:       intr.NewSensorWithPlatform("mqtt", cfg),
		Availability: newMqttAvailability(cfg),
		stateTopic:   cfg["state_topic"],
	}

	if templ := cfg["value_template"]; len(templ) > 0 {
		t, err := template.New("valueTemplate").Parse(templ)

		if err == nil {
			ms.valueTemplate = t
		}
	}

	return ms
}
