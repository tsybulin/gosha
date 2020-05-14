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

func (hk *hksecsy) targetStateRemoteUpdate(state int) {
	srv := svc.NewRegistry(hk.eventBus).GetDomainService(cmp.DomainAlarm).(svc.Alarm)
	cs := hk.accsecsy.SecuritySystem.SecuritySystemCurrentState

	switch state {
	case characteristic.SecuritySystemTargetStateDisarm:
		srv.Disarm(hk.alarm.GetID())
		cs.SetValue(characteristic.SecuritySystemCurrentStateDisarmed)
	case characteristic.SecuritySystemTargetStateStayArm:
		srv.ArmHome(hk.alarm.GetID())
		cs.SetValue(characteristic.SecuritySystemCurrentStateStayArm)
	case characteristic.SecuritySystemTargetStateAwayArm:
		srv.ArmAway(hk.alarm.GetID())
		cs.SetValue(characteristic.SecuritySystemCurrentStateAwayArm)
	case characteristic.SecuritySystemTargetStateNightArm:
		srv.ArmNight(hk.alarm.GetID())
		cs.SetValue(characteristic.SecuritySystemCurrentStateNightArm)
	}

}
func (hk *hksecsy) currentStateRemoteUpdate(state int) {
	srv := svc.NewRegistry(hk.eventBus).GetDomainService(cmp.DomainAlarm).(svc.Alarm)

	switch state {
	case characteristic.SecuritySystemCurrentStateDisarmed:
		srv.Disarm(hk.alarm.GetID())
	case characteristic.SecuritySystemCurrentStateStayArm:
		srv.ArmHome(hk.alarm.GetID())
	case characteristic.SecuritySystemCurrentStateAwayArm:
		srv.ArmAway(hk.alarm.GetID())
	case characteristic.SecuritySystemCurrentStateNightArm:
		srv.ArmNight(hk.alarm.GetID())
	}

}

func (hk *hksecsy) stateLocalUpdate(state string) {
	cs := hk.accsecsy.SecuritySystem.SecuritySystemCurrentState
	ts := hk.accsecsy.SecuritySystem.SecuritySystemTargetState

	switch state {
	case cmp.AlarmStateDisarmed.String():
		ts.SetValue(characteristic.SecuritySystemTargetStateDisarm)
		cs.SetValue(characteristic.SecuritySystemCurrentStateDisarmed)
	case cmp.AlarmStateArmedHome.String():
		ts.SetValue(characteristic.SecuritySystemTargetStateStayArm)
		cs.SetValue(characteristic.SecuritySystemCurrentStateStayArm)
	case cmp.AlarmStateArmedAway.String():
		ts.SetValue(characteristic.SecuritySystemTargetStateAwayArm)
		cs.SetValue(characteristic.SecuritySystemCurrentStateAwayArm)
	case cmp.AlarmStateArmedNight.String():
		ts.SetValue(characteristic.SecuritySystemTargetStateNightArm)
		cs.SetValue(characteristic.SecuritySystemCurrentStateNightArm)
	case cmp.AlarmStateTriggered.String():
		cs.SetValue(characteristic.SecuritySystemCurrentStateAlarmTriggered)
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

	acc.SecuritySystem.SecuritySystemCurrentState.OnValueRemoteUpdate(hks.currentStateRemoteUpdate)
	acc.SecuritySystem.SecuritySystemTargetState.OnValueRemoteUpdate(hks.targetStateRemoteUpdate)

	return hks
}
