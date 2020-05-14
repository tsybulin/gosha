package ptf

import (
	"context"
	"strconv"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/logger"
)

type homekit struct {
	Platform
	name       string
	pin        string
	includes   map[string]int
	components map[string]cmp.Component
	accs       map[string]interface{}
	eventBus   evt.Bus
}

func (hk *homekit) registerComponent(c cmp.Component) {
	if _, ok := hk.includes[c.GetID()]; ok {
		hk.components[c.GetID()] = c
	}
}

func (hk *homekit) componentStateChanged(event evt.Message) {
	if "event" != event.Type || "state_changed" != event.Event.EventType || event.Event.Data.NewState == nil {
		return
	}

	acc := hk.accs[event.Event.Data.EntityID]

	if acc == nil {
		return
	}

	switch acc.(type) {
	case *hklight:
		acc.(*hklight).brightnessLocalUpdate(int(event.Event.Data.NewState.Attributes["brightness"].(int16)))
		acc.(*hklight).stateLocalUpdate(event.Event.Data.NewState.State)
	case *hkswitch:
		acc.(*hkswitch).stateLocalUpdate(event.Event.Data.NewState.State)
	case *hksecsy:
		acc.(*hksecsy).stateLocalUpdate(event.Event.Data.NewState.State)
	}
}

func (hk *homekit) componentRegistrationFinished(b bool) {
	bridge := accessory.NewBridge(accessory.Info{
		Name:             hk.name,
		Manufacturer:     "Pavel Tsybulin",
		Model:            "Gosha",
		ID:               1,
		SerialNumber:     "homekit.bridge",
		FirmwareRevision: "0.1.0",
	})

	accessories := make([]*accessory.Accessory, 0)

	for id, c := range hk.components {
		if c.GetDomain() == cmp.DomainLight {
			acc := hk.newhklight(c.(cmp.Light))
			accessories = append(accessories, acc.acclight.Accessory)
			hk.accs[id] = acc
		} else if c.GetDomain() == cmp.DomainSwitch {
			acc := hk.newhkswitch(c.(cmp.Switch))
			accessories = append(accessories, acc.accswitch.Accessory)
			hk.accs[id] = acc
		} else if c.GetDomain() == cmp.DomainAlarm {
			acc := hk.newhksecsy(c.(cmp.Alarm))
			accessories = append(accessories, acc.accsecsy.Accessory)
			hk.accs[id] = acc
		}
	}

	transport, err := hc.NewIPTransport(
		hc.Config{
			Pin:         hk.pin,
			SetupId:     "GOSH",
			StoragePath: "./config/homekit",
		}, bridge.Accessory, accessories[0:]...)

	if err != nil {
		hk.eventBus.Publish(logger.Topic, logger.LevelError, "HomeKit.transport error", err.Error())
	} else {
		hc.OnTermination(func() {
			transport.Stop()
		})

		go func() {
			ctx := context.Background()
			transport.Start()
			<-ctx.Done()
		}()

	}
}

func (hk *homekit) Start(eventBus evt.Bus) {
	hk.eventBus = eventBus

	hk.eventBus.Subscribe(cmp.RegisterTopic, "HomeKit.registerComponent", hk.registerComponent)
	hk.eventBus.Subscribe(evt.StateChangedTopic, "HomeKit.componentStateChanged", hk.componentStateChanged)
	hk.eventBus.SubscribeOnceAsync(cmp.RegistrationFinishedTopic, "HomeKit.componentRegistrationFinished", hk.componentRegistrationFinished)

	hk.eventBus.Publish(ReadyTopic, cmp.DomainHomeKit)
}

// NewHomeKitPlatform ...
func NewHomeKitPlatform(cfg map[string]string) Platform {
	hk := &homekit{
		Platform:   newPlatform("homekit"),
		name:       cfg["homekit"],
		pin:        cfg["pin"],
		includes:   make(map[string]int),
		components: make(map[string]cmp.Component),
		accs:       make(map[string]interface{}),
	}

	for k, v := range cfg {
		if k == "homekit" || k == "pin" {
			continue
		}

		if d, err := strconv.Atoi(v); err == nil {
			hk.includes[k] = d
		}

	}

	return hk
}
