package intr

import (
	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
)

type light struct {
	cmp.Switch
	brightness int16
}

// GetBrightness ...
func (l *light) GetBrightness() int16 {
	return l.brightness
}

// SetBrightness ...
func (l *light) SetBrightness(brightness int16) {
	l.brightness = brightness
	// log.Println("--- brightness", l.GetID(), l.brightness)
}

func (l *light) GetState() evt.State {
	state := l.Switch.GetState()

	state.Attributes["supported_features"] = 1
	state.Attributes["brightness"] = l.GetBrightness()

	return state
}

func NewLight(cfg map[string]string) cmp.Light {
	return &light{
		Switch:     NewDomainSwitchWithPlatform(cmp.DomainLight, cfg["light"], "internal"),
		brightness: 0,
	}
}

func NewLightWithPlatform(platform string, cfg map[string]string) cmp.Light {
	return &light{
		Switch:     NewDomainSwitchWithPlatform(cmp.DomainLight, cfg["light"], platform),
		brightness: 0,
	}
}
