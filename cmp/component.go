package cmp

import "github.com/tsybulin/gosha/evt"

// Statefull ...
type Statefull interface {
	GetID() string
	GetState() evt.State
}

// Component ...
type Component interface {
	Statefull
	GetID() string
	GetPlatform() string
	GetDomain() Domain
}
