package cmp

// Domain ...
type Domain int8

const (
	DomainSwitch Domain = iota
	DomainTimer
	DomainLight
	DomainSensor
	DomainBinarySensor
	DomainGroup
	DomainMqtt
	DomainScript
	DomainWeather
	DomainTelegram
	DomainAlarm
	DomainHomeKit
	DomainCountdown
)

// Domains ...
var Domains = []Domain{
	DomainSwitch,
	DomainTimer,
	DomainLight,
	DomainSensor,
	DomainBinarySensor,
	DomainBinarySensor,
	DomainGroup,
	DomainMqtt,
	DomainScript,
	DomainWeather,
	DomainTelegram,
	DomainAlarm,
	DomainHomeKit,
	DomainCountdown,
}

// DomainNames ...
var DomainNames = []string{
	"switch",
	"timer",
	"light",
	"sensor",
	"binary_sensor",
	"group",
	"mqtt",
	"script",
	"weather",
	"telegram",
	"alarm",
	"homekit",
	"countdown",
}

func (d Domain) String() string {
	return DomainNames[d]
}
