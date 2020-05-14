package mqtt

import (
	"html/template"

	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/cmp/intr"
	"github.com/tsybulin/gosha/evt"
)

// Light ...
type Light interface {
	cmp.Light
	Availability
	GetStateTopic() string
	GetPayloadOn() string
	GetPayloadOff() string
	GetCommandTopic() string
	GetStateValueTemplate() *template.Template
	GetBrightnessStateTopic() string
	GetBrightnessValueTemplate() *template.Template
	GetBrightnessCommandTopic() string
}

type mqttLight struct {
	cmp.Light
	Availability
	stateTopic              string
	payloadOn               string
	payloadOff              string
	commandTopic            string
	stateValueTemplate      *template.Template
	brightnessStateTopic    string
	brightnessValueTemplate *template.Template
	brightnessCommandTopic  string
}

func (ml *mqttLight) GetStateTopic() string {
	return ml.stateTopic
}

func (ml *mqttLight) GetPayloadOn() string {
	return ml.payloadOn
}

func (ml *mqttLight) GetPayloadOff() string {
	return ml.payloadOff
}

func (ml *mqttLight) GetCommandTopic() string {
	return ml.commandTopic
}

func (ml *mqttLight) GetStateValueTemplate() *template.Template {
	return ml.stateValueTemplate
}

func (ml *mqttLight) GetBrightnessStateTopic() string {
	return ml.brightnessStateTopic
}

func (ml *mqttLight) GetBrightnessValueTemplate() *template.Template {
	return ml.brightnessValueTemplate
}

func (ml *mqttLight) GetBrightnessCommandTopic() string {
	return ml.brightnessCommandTopic
}

func (ml *mqttLight) GetState() evt.State {
	state := ml.Light.GetState()
	state.Attributes["available"] = ml.GetAvailable()
	return state
}

// NewMqttLight ...
func NewMqttLight(cfg map[string]string) Light {
	ml := &mqttLight{
		Light:                  intr.NewLightWithPlatform("mqtt", cfg),
		Availability:           newMqttAvailability(cfg),
		stateTopic:             cfg["state_topic"],
		payloadOn:              cfg["payload_on"],
		payloadOff:             cfg["payload_off"],
		commandTopic:           cfg["command_topic"],
		brightnessStateTopic:   cfg["brightness_state_topic"],
		brightnessCommandTopic: cfg["brightness_command_topic"],
	}

	if templ := cfg["state_value_template"]; len(templ) > 0 {
		t, err := template.New("stateValueTemplate").Parse(templ)

		if err == nil {
			ml.stateValueTemplate = t
		}
	}

	if templ := cfg["brightness_value_template"]; len(templ) > 0 {
		t, err := template.New("brightnessValueTemplate").Parse(templ)

		if err == nil {
			ml.brightnessValueTemplate = t
		}
	}

	return ml
}
