package cmp

// AlarmState ...
type AlarmState int

const (
	AlarmStateUnknown AlarmState = iota
	AlarmStateDisarmed
	AlarmStateArmedHome
	AlarmStateArmedAway
	AlarmStateArmedNight
	AlarmStateTriggered
)

func (s AlarmState) String() string {
	return [...]string{
		"unknown",
		"disarmed",
		"armed_home",
		"armed_away",
		"armed_night",
		"triggered",
	}[s]
}

// Alarm ...
type Alarm interface {
	Component
	AlarmState() AlarmState
	SetAlarmState(AlarmState)
}
