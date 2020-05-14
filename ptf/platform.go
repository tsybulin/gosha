package ptf

import (
	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
)

const (
	// RegisterTopic ...
	RegisterTopic = "platform:register"

	// ReadyTopic ...
	ReadyTopic = "platform:ready"
)

// Platform ...
type Platform interface {
	GetPlatform() string
	Start(evt.Bus)
	Push(cmp.Component, string)
}

type platform struct {
	platform string
}

func (p *platform) GetPlatform() string {
	return p.platform
}

func (p *platform) Push(cmp.Component, string) {
	// nothing to do here
}

func (p *platform) Start(evt.Bus) {

}

func newPlatform(p string) Platform {
	return &platform{
		platform: p,
	}
}

func eventFor(c cmp.Component, fn func()) evt.Message {
	oldState := c.GetState()
	fn()
	newState := c.GetState()

	message := evt.Message{
		ID:   0,
		Type: "event",
		Event: &evt.Event{
			EventType: "state_changed",
			Data: evt.Data{
				EntityID: c.GetID(),
				OldState: &oldState,
				NewState: &newState,
			},
		},
	}

	return message
}
