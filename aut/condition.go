package aut

// Condition ...
type Condition interface {
	GetPlatform() string
	Satisfied() bool
}

// StateCondition ...
type StateCondition interface {
	Condition
	GetEntityID() string
	GetState() string
	SatisfiedState(string, string) bool
}
