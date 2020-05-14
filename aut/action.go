package aut

// Action ...
type Action interface {
	GetService() string
	GetAction() string
	GetComponent() string
	Execute()
}
