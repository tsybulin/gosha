package cmp

// Group ...
type Group interface {
	Component
	AddComponent(Component) bool
	GetComponents() []Component
	GetComponent(string) Component
	IsOn() bool
}
