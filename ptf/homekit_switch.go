package ptf

import (
	"strings"

	"github.com/brutella/hc/accessory"
	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/svc"
)

type hkswitch struct {
	accswitch *accessory.Switch
	_switch   cmp.Switch
	eventBus  evt.Bus
}

func (s *hkswitch) onRemoteUpdate(on bool) {
	if on {
		svc.NewRegistry(s.eventBus).GetDomainService(cmp.DomainSwitch).(svc.Switch).TurnOn(s._switch.GetID())
	} else {
		svc.NewRegistry(s.eventBus).GetDomainService(cmp.DomainSwitch).(svc.Switch).TurnOff(s._switch.GetID())
	}
}

func (s *hkswitch) stateLocalUpdate(state string) {
	s.accswitch.Switch.On.SetValue(state == "on")
}

func (hk *homekit) newhkswitch(c cmp.Switch) *hkswitch {
	hkname := strings.ReplaceAll(c.GetID(), "switch.", "")
	hkname = strings.ReplaceAll(hkname, "_", " ")

	acc := accessory.NewSwitch(accessory.Info{
		Name:             hkname,
		ID:               uint64(hk.includes[c.GetID()]),
		SerialNumber:     c.GetID(),
		Model:            c.GetDomain().String(),
		Manufacturer:     "Gosha",
		FirmwareRevision: "0.1.0",
	})

	hks := &hkswitch{
		accswitch: acc,
		_switch:   c,
		eventBus:  hk.eventBus,
	}

	acc.Switch.On.SetValue(hks._switch.IsOn())
	acc.Switch.On.OnValueRemoteUpdate(hks.onRemoteUpdate)

	return hks
}
