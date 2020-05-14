package cmp

// Light ...
type Light interface {
	Switch
	GetBrightness() int16
	SetBrightness(int16)
}
