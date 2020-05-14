package aut

const RegisterTopic = "automation:register"

// Automation ...
type Automation interface {
	GetID() string
	GetTriggers() []Trigger
	GetContitions() []Condition
	GetActions() []Action
	IsBusy() bool
	SetBusy(bool)
}
