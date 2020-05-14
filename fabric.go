package gosha

import (
	"github.com/tsybulin/gosha/aut"
	"github.com/tsybulin/gosha/aut/ai"
	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/cmp/intr"
	"github.com/tsybulin/gosha/cmp/mqtt"
)

// CreateComponent create new Component from config
func CreateComponent(cfg map[string]string) cmp.Component {
	if id := cfg["switch"]; len(id) > 0 {
		platform := cfg["platform"]
		if len(platform) == 0 {
			return intr.NewInternalSwitch(id)

		} else if platform == "mqtt" {
			return mqtt.NewMqttSwitch(cfg)
		}
	}

	if id := cfg["timer"]; len(id) > 0 {
		return intr.NewTimer(cfg)
	}

	if id := cfg["light"]; len(id) > 0 {
		platform := cfg["platform"]
		if len(platform) == 0 {
			return intr.NewLight(cfg)
		} else if platform == "mqtt" {
			return mqtt.NewMqttLight(cfg)
		}
	}

	if id := cfg["binary_sensor"]; len(id) > 0 {
		platform := cfg["platform"]
		if platform == "mqtt" {
			return mqtt.NewMqttBinarySensor(cfg)
		}
	}

	if id := cfg["sensor"]; len(id) > 0 {
		platform := cfg["platform"]
		if platform == "mqtt" {
			return mqtt.NewMqttSensor(cfg)
		}

		if platform == "weather" {
			return intr.NewWeather(cfg)
		}
	}

	if id := cfg["telegram"]; len(id) > 0 {
		return intr.NewTelegram(cfg)
	}

	if id := cfg["alarm"]; len(id) > 0 {
		platform := cfg["platform"]
		if len(platform) == 0 {
			return intr.NewAlarm(cfg)
		} else if platform == "mqtt" {
			return mqtt.NewMqttAlarm(cfg)
		}

	}

	return nil
}

// CreateGroup ...
func CreateGroup(id string, es []string) cmp.Group {
	return intr.NewGroup(id, es)
}

// NewAutomations ...
func NewAutomations(acfg []struct {
	Automation string
	Triggers   []map[string]string
	Conditions []map[string]string
	Actions    []map[string]string
}) []aut.Automation {
	return ai.NewAutomations(acfg)
}

// NewScripts ...
func NewScripts(scfg []struct {
	Script  string
	Actions []map[string]string
}) []aut.Script {
	return ai.NewScripts(scfg)
}
