package cmp

// Switchable ...
type Switchable interface {
	IsOn() bool
	TurnOn()
	TurnOff()
	Toggle()
	GetOnString() string
}

// Switch ...
type Switch interface {
	Component
	Switchable
}
