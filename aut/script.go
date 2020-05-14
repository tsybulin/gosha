package aut

import "github.com/tsybulin/gosha/evt"

const ScriptRegisterTopic = "script:register"

// Script ...
type Script interface {
	GetID() string
	Execute()
	IsActive() bool
	SetActive(bool)
	GetState() evt.State
}
