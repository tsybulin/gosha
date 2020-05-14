package svc

import (
	"github.com/tsybulin/gosha/cmp"
	"github.com/tsybulin/gosha/evt"
)

// Method ...
type Method struct {
	Method     string   `json:"method"`
	Parameters []string `json:"parameters,omitempty"`
}

// Description ...
type Description struct {
	ID      string   `json:"id"`
	Methods []Method `json:"methods,omitempty"`
}

// StateResult ...
type StateResult struct {
	ID      int         `json:"id"`
	Type    string      `json:"type"`
	Success bool        `json:"success"`
	Result  []evt.State `json:"result"`
}

// Service ...
type Service interface {
	GetID() string
	GetDomain() cmp.Domain
	Description() Description
	Components() []string
	States() StateResult
}

type service struct {
	id     string
	domain cmp.Domain
}

func (s *service) GetID() string {
	return s.id
}

func (s *service) GetDomain() cmp.Domain {
	return s.domain
}

func (s *service) Description() Description {
	return Description{
		ID: s.GetID(),
	}
}

func (s *service) Components() []string {
	components := make([]string, 0)
	return components
}

func (s *service) States() StateResult {
	return StateResult{
		ID:      1,
		Type:    "result",
		Success: true,
	}
}

func stateName(b bool) string {
	if b {
		return "on"
	} else {
		return "off"
	}
}

func newService(domain cmp.Domain) Service {
	return &service{
		domain: domain,
		id:     "service." + domain.String(),
	}
}

func eventFor(c cmp.Statefull, fn func()) evt.Message {
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
