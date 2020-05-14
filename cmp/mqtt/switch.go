package mqtt

import (
	"html/template"

	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/cmp/intr"
	"github.com/tsybulin/gosha/evt"
)

// Switch ...
type Switch interface {
	cmp.Switch
	Availability
	GetStateTopic() string
	GetPayloadOn() string
	GetPayloadOff() string
	GetCommandTopic() string
	GetStateValueTemplate() *template.Template
}

type mqttSwitch struct {
	cmp.Switch
	Availability
	stateTopic         string
	payloadOn          string
	payloadOff         string
	commandTopic       string
	stateValueTemplate *template.Template
}

func (ms *mqttSwitch) GetStateTopic() string {
	return ms.stateTopic
}

func (ms *mqttSwitch) GetPayloadOn() string {
	return ms.payloadOn
}

func (ms *mqttSwitch) GetPayloadOff() string {
	return ms.payloadOff
}

func (ms *mqttSwitch) GetCommandTopic() string {
	return ms.commandTopic
}

func (ms *mqttSwitch) GetStateValueTemplate() *template.Template {
	return ms.stateValueTemplate
}

func (ms *mqttSwitch) GetState() evt.State {
	state := ms.Switch.GetState()
	state.Attributes["available"] = ms.GetAvailable()
	return state
}

// NewMqttSwitch ...
func NewMqttSwitch(cfg map[string]string) Switch {
	sw := &mqttSwitch{
		Switch:       intr.NewSwitchWithPlatform(cfg["switch"], cfg["platform"]),
		Availability: newMqttAvailability(cfg),
		stateTopic:   cfg["state_topic"],
		payloadOn:    cfg["payload_on"],
		payloadOff:   cfg["payload_off"],
		commandTopic: cfg["command_topic"],
	}

	if templ := cfg["state_value_template"]; len(templ) > 0 {
		t, err := template.New("stateValueTemplate").Parse(templ)

		if err == nil {
			sw.stateValueTemplate = t
		}
	}

	return sw
}
