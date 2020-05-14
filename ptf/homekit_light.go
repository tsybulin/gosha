package ptf

import (
	"strings"

	"github.com/brutella/hc/accessory"
	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/svc"
)

type hklight struct {
	acclight *accessory.ColoredLightbulb
	light    cmp.Light
	eventBus evt.Bus
}

func (s *hklight) onRemoteUpdate(on bool) {
	if on {
		svc.NewRegistry(s.eventBus).GetDomainService(cmp.DomainLight).(svc.Light).TurnOn(s.light.GetID())
	} else {
		svc.NewRegistry(s.eventBus).GetDomainService(cmp.DomainLight).(svc.Light).TurnOff(s.light.GetID())
	}
}

func (s *hklight) brightnessRemoteUpdate(b int) {
	svc.NewRegistry(s.eventBus).GetDomainService(cmp.DomainLight).(svc.Light).SetBrightness(s.light.GetID(), int16(b))
}

func (s *hklight) stateLocalUpdate(state string) {
	s.acclight.Lightbulb.On.SetValue(state == "on")
}

func (s *hklight) brightnessLocalUpdate(b int) {
	s.acclight.Lightbulb.Brightness.SetValue(b)
}

func (hk *homekit) newhklight(c cmp.Light) *hklight {
	hkname := strings.ReplaceAll(c.GetID(), "light.", "")
	hkname = strings.ReplaceAll(hkname, "_", " ")

	acc := accessory.NewColoredLightbulb(accessory.Info{
		Name:             hkname,
		ID:               uint64(hk.includes[c.GetID()]),
		SerialNumber:     c.GetID(),
		Model:            c.GetDomain().String(),
		Manufacturer:     "Gosha",
		FirmwareRevision: "0.1.0",
	})

	hkl := &hklight{
		acclight: acc,
		light:    c,
		eventBus: hk.eventBus,
	}

	acc.Lightbulb.On.SetValue(hkl.light.IsOn())
	acc.Lightbulb.Brightness.SetValue(int(hkl.light.GetBrightness()))
	acc.Lightbulb.On.OnValueRemoteUpdate(hkl.onRemoteUpdate)
	acc.Lightbulb.Brightness.OnValueRemoteUpdate(hkl.brightnessRemoteUpdate)

	return hkl
}
