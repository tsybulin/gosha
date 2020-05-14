package intr

import (
	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
)

type internalSwitch struct {
	cmp.Component
	cmp.Switchable
	power bool
}

func (s *internalSwitch) IsOn() bool {
	return s.power
}

func (s *internalSwitch) GetOnString() string {
	if s.IsOn() {
		return "on"
	} else {
		return "off"
	}
}

func (s *internalSwitch) TurnOn() {
	s.power = true
	// log.Println("--- on power", s.GetID(), s.power)
}

func (s *internalSwitch) TurnOff() {
	s.power = false
	// log.Println("--- off power", s.GetID(), s.power)
}

func (s *internalSwitch) Toggle() {
	s.power = !s.power
	// log.Println("--- toggle power", s.GetID(), s.power)
}

func (s *internalSwitch) GetState() evt.State {
	return evt.State{
		EntityID:   s.GetID(),
		State:      s.GetOnString(),
		Attributes: make(map[string]interface{}, 0),
	}
}

func NewInternalSwitch(id string) cmp.Switch {
	return &internalSwitch{
		Component: NewComponent(cmp.DomainSwitch, id, "internal"),
		power:     false,
	}
}

func NewSwitchWithPlatform(id, platform string) cmp.Switch {
	return &internalSwitch{
		Component: NewComponent(cmp.DomainSwitch, id, platform),
		power:     false,
	}
}

func NewDomainSwitchWithPlatform(domain cmp.Domain, id, platform string) cmp.Switch {
	return &internalSwitch{
		Component: NewComponent(domain, id, platform),
		power:     false,
	}
}
