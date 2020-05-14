package ptf

import (
	"strings"

	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
	"github.com/tsybulin/gosha/svc"
)

const (
	SecuritySystemAlarmTypeNoAlarm = 0
	SecuritySystemAlarmTypeAlarm   = 1
)

type secsy struct {
	*accessory.Accessory
	AlarmType      *characteristic.SecuritySystemAlarmType
	SecuritySystem *service.SecuritySystem
}

func newsecsy(info accessory.Info) *secsy {
	acc := secsy{}
	acc.Accessory = accessory.New(info, accessory.TypeSecuritySystem)
	acc.SecuritySystem = service.NewSecuritySystem()

	acc.AlarmType = characteristic.NewSecuritySystemAlarmType()
	acc.AlarmType.SetValue(SecuritySystemAlarmTypeAlarm)
	acc.SecuritySystem.AddCharacteristic(acc.AlarmType.Characteristic)
	acc.AddService(acc.SecuritySystem.Service)

	return &acc
}

type hksecsy struct {
	accsecsy *secsy
	alarm    cmp.Alarm
	eventBus evt.Bus
}

func (hk *hksecsy) stateRemoteUpdate(state int) {
	switch state {
	case characteristic.SecuritySystemTargetStateDisarm:
		svc.NewRegistry(hk.eventBus).GetDomainService(cmp.DomainAlarm).(svc.Alarm).Disarm(hk.alarm.GetID())
		hk.accsecsy.SecuritySystem.SecuritySystemTargetState.SetValue(characteristic.SecuritySystemTargetStateDisarm)
		hk.accsecsy.SecuritySystem.SecuritySystemCurrentState.SetValue(characteristic.SecuritySystemCurrentStateDisarmed)
	case characteristic.SecuritySystemTargetStateStayArm:
		svc.NewRegistry(hk.eventBus).GetDomainService(cmp.DomainAlarm).(svc.Alarm).ArmHome(hk.alarm.GetID())
		hk.accsecsy.SecuritySystem.SecuritySystemTargetState.SetValue(characteristic.SecuritySystemTargetStateStayArm)
		hk.accsecsy.SecuritySystem.SecuritySystemCurrentState.SetValue(characteristic.SecuritySystemCurrentStateStayArm)
	case characteristic.SecuritySystemTargetStateAwayArm:
		svc.NewRegistry(hk.eventBus).GetDomainService(cmp.DomainAlarm).(svc.Alarm).ArmAway(hk.alarm.GetID())
		hk.accsecsy.SecuritySystem.SecuritySystemTargetState.SetValue(characteristic.SecuritySystemTargetStateAwayArm)
		hk.accsecsy.SecuritySystem.SecuritySystemCurrentState.SetValue(characteristic.SecuritySystemCurrentStateAwayArm)
	case characteristic.SecuritySystemTargetStateNightArm:
		svc.NewRegistry(hk.eventBus).GetDomainService(cmp.DomainAlarm).(svc.Alarm).ArmNight(hk.alarm.GetID())
		hk.accsecsy.SecuritySystem.SecuritySystemTargetState.SetValue(characteristic.SecuritySystemTargetStateNightArm)
		hk.accsecsy.SecuritySystem.SecuritySystemCurrentState.SetValue(characteristic.SecuritySystemCurrentStateNightArm)
	}

}

func (hk *hksecsy) stateLocalUpdate(state string) {
	switch state {
	case cmp.AlarmStateDisarmed.String():
		hk.accsecsy.SecuritySystem.SecuritySystemCurrentState.SetValue(characteristic.SecuritySystemCurrentStateDisarmed)
	case cmp.AlarmStateArmedHome.String():
		hk.accsecsy.SecuritySystem.SecuritySystemCurrentState.SetValue(characteristic.SecuritySystemCurrentStateStayArm)
	case cmp.AlarmStateArmedAway.String():
		hk.accsecsy.SecuritySystem.SecuritySystemCurrentState.SetValue(characteristic.SecuritySystemCurrentStateAwayArm)
	case cmp.AlarmStateArmedNight.String():
		hk.accsecsy.SecuritySystem.SecuritySystemCurrentState.SetValue(characteristic.SecuritySystemCurrentStateNightArm)
	case cmp.AlarmStateTriggered.String():
		hk.accsecsy.SecuritySystem.SecuritySystemCurrentState.SetValue(characteristic.SecuritySystemCurrentStateAlarmTriggered)
	}
}

func (hk *homekit) newhksecsy(c cmp.Alarm) *hksecsy {
	hkname := strings.ReplaceAll(c.GetID(), "alarm.", "")
	hkname = strings.ReplaceAll(hkname, "_", " ")

	acc := newsecsy(accessory.Info{
		Name:             hkname,
		ID:               uint64(hk.includes[c.GetID()]),
		SerialNumber:     c.GetID(),
		Model:            c.GetDomain().String(),
		Manufacturer:     "Gosha",
		FirmwareRevision: "0.1.0",
	})

	hks := &hksecsy{
		accsecsy: acc,
		alarm:    c,
		eventBus: hk.eventBus,
	}

	switch c.AlarmState() {
	case cmp.AlarmStateDisarmed:
		acc.SecuritySystem.SecuritySystemCurrentState.SetValue(characteristic.SecuritySystemCurrentStateDisarmed)
		acc.SecuritySystem.SecuritySystemTargetState.SetValue(characteristic.SecuritySystemTargetStateDisarm)
	case cmp.AlarmStateArmedHome:
		acc.SecuritySystem.SecuritySystemCurrentState.SetValue(characteristic.SecuritySystemCurrentStateStayArm)
		acc.SecuritySystem.SecuritySystemTargetState.SetValue(characteristic.SecuritySystemTargetStateStayArm)
	case cmp.AlarmStateArmedAway:
		acc.SecuritySystem.SecuritySystemCurrentState.SetValue(characteristic.SecuritySystemCurrentStateAwayArm)
		acc.SecuritySystem.SecuritySystemTargetState.SetValue(characteristic.SecuritySystemTargetStateAwayArm)
	case cmp.AlarmStateArmedNight:
		acc.SecuritySystem.SecuritySystemCurrentState.SetValue(characteristic.SecuritySystemCurrentStateNightArm)
		acc.SecuritySystem.SecuritySystemTargetState.SetValue(characteristic.SecuritySystemTargetStateNightArm)
	}

	acc.SecuritySystem.SecuritySystemTargetState.OnValueRemoteUpdate(hks.stateRemoteUpdate)

	return hks
}
