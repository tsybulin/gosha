package cmp

// Sensor ...
type Sensor interface {
	Component
	GetValue() interface{}
	SetValue(interface{})
	GetDeviceClass() string
	GetUnitOfMeasurement() string
}

// BinarySensor ...
type BinarySensor interface {
	Component
	IsOn() bool
	SetOn(bool)
	GetDeviceClass() string
	GetOnString() string
}
