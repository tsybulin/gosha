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

// CompareCondition ...
type CompareCondition interface {
	Condition
	GetEntityID() string
	GetValue() float64
	SatisfiedCompare(string, float64) bool
}
